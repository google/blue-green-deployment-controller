package bluegreendeployment

import (
	"fmt"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/bgd/pkg/apis/controller/v1alpha1"
	controllerv1alpha1 "k8s.io/bgd/pkg/apis/controller/v1alpha1"
	controllerv1alpha1client "k8s.io/bgd/pkg/client/clientset_generated/clientset/typed/controller/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// int64Ptr returns a pointer of int64.
func int64Ptr(i int64) *int64 {
	return &i
}

// newReplicaSet returns a new ReplicaSet struct based on .spec field of BlueGreenDeployment object.
// If the ReplicaSet has different color label than .status.ActiveReplicaSetColor field of the
// BlueGreenDeployment object, it is considered inactive and is specified to have zero replica.
func newReplicaSet(b *v1alpha1.BlueGreenDeployment, color string) *appsv1.ReplicaSet {
	replicas := b.Spec.Replicas
	if b.Status.ActiveReplicaSetColor != color {
		replicas = int32(0)
	}

	rsName := BlueRSName
	if color != BlueColor {
		rsName = GreenRSName
	}

	return &appsv1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ReplicaSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rsName,
			Namespace: b.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(b, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    BlueGreenDeploymentKind,
				}),
			},
		},
		Spec: appsv1.ReplicaSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"color": color},
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"color": color},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:                     b.Spec.PodSpec.Containers[0].Name,
							Image:                    b.Spec.PodSpec.Containers[0].Image,
							Command:                  b.Spec.PodSpec.Containers[0].Command,
							TerminationMessagePath:   b.Spec.PodSpec.Containers[0].TerminationMessagePath,
							TerminationMessagePolicy: b.Spec.PodSpec.Containers[0].TerminationMessagePolicy,
							ImagePullPolicy:          b.Spec.PodSpec.Containers[0].ImagePullPolicy,
						},
					},
					TerminationGracePeriodSeconds: b.Spec.PodSpec.TerminationGracePeriodSeconds,
					DNSPolicy:                     b.Spec.PodSpec.DNSPolicy,
					SecurityContext:               b.Spec.PodSpec.SecurityContext,
					SchedulerName:                 b.Spec.PodSpec.SchedulerName,
				},
			},
		},
	}
}

// newService returns a new service struct pointing to blue ReplicaSet by default.
func newService(b *v1alpha1.BlueGreenDeployment) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "core/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceName,
			Namespace: b.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{LabelColorKey: BlueColor},
			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       80,
					TargetPort: intstr.FromInt(443),
				},
			},
		},
	}
}

// createReplicaSetsIfUnavailable creates blue and green ReplicaSets if they are unavailable.
func createReplicaSetsIfUnavailable(b *v1alpha1.BlueGreenDeployment, c *kubernetes.Clientset) (*appsv1.ReplicaSet, *appsv1.ReplicaSet, error) {
	blueRS, err := createReplicaSetIfUnavailable(b, c, BlueRSName, BlueColor)
	if err != nil {
		return nil, nil, fmt.Errorf("when creating blue ReplicaSet, %v", err)
	}

	greenRS, err := createReplicaSetIfUnavailable(b, c, GreenRSName, GreenColor)
	if err != nil {
		return nil, nil, fmt.Errorf("when creating green ReplicaSet, %v", err)
	}

	return blueRS, greenRS, nil
}

// createReplicaSetsIfUnavailable creates a ReplicaSet with given color label if it is unavailable.
func createReplicaSetIfUnavailable(b *v1alpha1.BlueGreenDeployment, c *kubernetes.Clientset, rsName, rsColor string) (*appsv1.ReplicaSet, error) {
	ns := b.Namespace
	rs, err := c.AppsV1().ReplicaSets(ns).Get(rsName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			rs, err = c.AppsV1().ReplicaSets(ns).Create(newReplicaSet(b, rsColor))
			if err != nil {
				return nil, fmt.Errorf("failed to create ReplicaSet %q: %v", rsName, err)
			}
			if err = waitAllReplicaSetPodsAvailable(c, rs); err != nil {
				// When not all pods become available within timeout, it may be due to invalid
				// container field input(s) such as invalid image name. We should not return an
				// error to ensure the ReplicaSet to be made inactive with valid container fields
				// is scaled down later down the sync call tree. Logging the failure as a message
				// should suffice as this controller does not support auto-rollback feature.
				log.Printf("some pods for ReplicaSet %q failed to become available: %v", rs.Name, err)
			}
			return rs, nil
		} else {
			return nil, fmt.Errorf("failed to get ReplicaSet %q: %v", rsName, err)
		}
	}

	return rs, nil
}

// updateBlueGreenDeployment updates a BlueGreenDeployment object based on given update function via polling.
func updateBlueGreenDeployment(c controllerv1alpha1client.ControllerV1alpha1Interface, bgdName, ns string, updateFunc func(*v1alpha1.BlueGreenDeployment)) (*v1alpha1.BlueGreenDeployment, error) {
	var blueGreenDeployment *v1alpha1.BlueGreenDeployment
	bgdClient := c.BlueGreenDeployments(ns)
	if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		bgd, err := bgdClient.Get(bgdName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get BlueGreenDeployment object: %v", err)
		}
		updateFunc(bgd)
		blueGreenDeployment, err = bgdClient.Update(bgd)
		if err != nil {
			return fmt.Errorf("failed to update BlueGreenDeployment object: %v", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return blueGreenDeployment, nil
}

// waitAllReplicaSetPodsAvailable waits all ReplicaSet pods to become available.
func waitAllReplicaSetPodsAvailable(c *kubernetes.Clientset, replicaSet *appsv1.ReplicaSet) error {
	desiredGeneration := replicaSet.Generation
	return wait.PollImmediate(PollInterval, PollTimeout, func() (bool, error) {
		rs, err := c.AppsV1().ReplicaSets(replicaSet.Namespace).Get(replicaSet.Name, metav1.GetOptions{})
		if err != nil {
			return false, fmt.Errorf("failed to get ReplicaSet %q: %v", replicaSet.Name, err)
		}
		return rs.Status.ObservedGeneration >= desiredGeneration && rs.Status.AvailableReplicas == *replicaSet.Spec.Replicas, nil
	})
}

// createServiceIfUnavailable creates service if it is unavailable.
func createServiceIfUnavailable(b *v1alpha1.BlueGreenDeployment, c *kubernetes.Clientset) (*corev1.Service, error) {
	ns := b.Namespace
	svc, err := c.CoreV1().Services(ns).Get(ServiceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			svc, err = c.CoreV1().Services(ns).Create(newService(b))
			if err != nil {
				return nil, fmt.Errorf("failed to create service %q: %v", ServiceName, err)
			}
		} else {
			return nil, fmt.Errorf("failed to get service %q: %v", ServiceName, err)
		}
	}

	return svc, nil
}

// scaleReplicaSet scales a ReplicaSet to given number of replicas. It updates .spec.replicas
// field of the ReplicaSet via polling and waits for new number of replicas to become available.
func scaleReplicaSet(c *kubernetes.Clientset, replicas int32, replicaSet *appsv1.ReplicaSet) error {
	replicaSet, err := updateReplicaSet(c, replicaSet.Name, replicaSet.Namespace, func(rs *appsv1.ReplicaSet) {
		*rs.Spec.Replicas = replicas
	})
	if err != nil {
		return fmt.Errorf("failed to update its .spec.replicas: %v", err)
	}

	if err = waitAllReplicaSetPodsAvailable(c, replicaSet); err != nil {
		log.Printf("some pods for ReplicaSet %q failed to become available: %v", replicaSet.Name, err)
	}
	return nil
}

// updateReplicaSet updates a ReplicaSet based on given update function via polling.
func updateReplicaSet(c *kubernetes.Clientset, rsName, ns string, updateFunc func(*appsv1.ReplicaSet)) (*appsv1.ReplicaSet, error) {
	var replicaSet *appsv1.ReplicaSet
	rsClient := c.AppsV1().ReplicaSets(ns)
	if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		rs, err := rsClient.Get(rsName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get ReplicaSet %q: %v", rsName, err)
		}
		updateFunc(rs)
		replicaSet, err = rsClient.Update(rs)
		if err != nil {
			return fmt.Errorf("failed to update ReplicaSet %q: %v", rsName, err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return replicaSet, nil
}

// replaceInactiveReplicaSet replaces the inactive ReplicaSet by first deleting it to make room for a
// new ReplicaSet. It then creates the new ReplicaSet based on pod spec of BlueGreenDeployment object.
func replaceInactiveReplicaSet(c *kubernetes.Clientset, b *v1alpha1.BlueGreenDeployment, inactiveRS *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error) {
	ns := b.Namespace
	rsName := inactiveRS.Name
	color := inactiveRS.Spec.Selector.MatchLabels[LabelColorKey]
	err := c.AppsV1().ReplicaSets(ns).Delete(rsName, &metav1.DeleteOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to delete inactive ReplicaSet %q: %v", rsName, err)
	}

	if err = waitReplicaSetCompletesDeletion(c, rsName, ns); err != nil {
		return nil, fmt.Errorf("inactive ReplicaSet %q failed to complete its deletion: %v", rsName, err)
	}

	replicaSet, err := createReplicaSetIfUnavailable(b, c, rsName, color)
	if err != nil {
		return nil, fmt.Errorf("when re-creating inactive ReplicaSet %q, %v", rsName, err)
	}

	return replicaSet, nil
}

// waitReplicaSetCompletesDeletion waits for a ReplicaSet to complete its deletion.
func waitReplicaSetCompletesDeletion(c *kubernetes.Clientset, rsName, ns string) error {
	return wait.PollImmediate(PollInterval, PollTimeout, func() (bool, error) {
		_, err := c.AppsV1().ReplicaSets(ns).Get(rsName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return true, nil
			}
			return false, fmt.Errorf("failed to get ReplicaSet %q: %v", rsName, err)
		}
		return false, nil
	})
}

// updateService updates a service based on given update function via polling.
func updateService(c *kubernetes.Clientset, ns string, updateFunc func(*corev1.Service)) (*corev1.Service, error) {
	var service *corev1.Service
	svcClient := c.CoreV1().Services(ns)
	if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		svc, err := svcClient.Get(ServiceName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get service %q: %v", ServiceName, err)
		}
		updateFunc(svc)
		service, err = svcClient.Update(svc)
		if err != nil {
			return fmt.Errorf("failed to update service %q: %v", ServiceName, err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return service, nil
}

// fillUnfilledBGDPodSpecFields fills unfilled pod spec fields of BlueGreenDeployment object with default values.
// WARNING: If user inputs an invalid value for a field, there is no validation to catch it. Corresponding
// ReplicaSet will not be created successfully, thus breaking underlying blue-green deployment mechanism.
func fillUnfilledBGDPodSpecFields(c controllerv1alpha1client.ControllerV1alpha1Interface, b *controllerv1alpha1.BlueGreenDeployment) (*controllerv1alpha1.BlueGreenDeployment, error) {
	var err error

	if b, err = updateBlueGreenDeployment(c, b.Name, b.Namespace, func(bgd *controllerv1alpha1.BlueGreenDeployment) {
		if bgd.Spec.PodSpec.Containers[0].TerminationMessagePath == "" {
			bgd.Spec.PodSpec.Containers[0].TerminationMessagePath = "/dev/termination-log"
		}
		if bgd.Spec.PodSpec.Containers[0].TerminationMessagePolicy == "" {
			bgd.Spec.PodSpec.Containers[0].TerminationMessagePolicy = "File"
		}
		if bgd.Spec.PodSpec.Containers[0].ImagePullPolicy == "" {
			bgd.Spec.PodSpec.Containers[0].ImagePullPolicy = "IfNotPresent"
		}
		if bgd.Spec.PodSpec.RestartPolicy == "" {
			bgd.Spec.PodSpec.RestartPolicy = "Always"
		}
		if bgd.Spec.PodSpec.TerminationGracePeriodSeconds == nil {
			bgd.Spec.PodSpec.TerminationGracePeriodSeconds = int64Ptr(0)
		}
		if bgd.Spec.PodSpec.DNSPolicy == "" {
			bgd.Spec.PodSpec.DNSPolicy = "ClusterFirst"
		}
		if bgd.Spec.PodSpec.SecurityContext == nil {
			bgd.Spec.PodSpec.SecurityContext = &corev1.PodSecurityContext{}
		}
		if bgd.Spec.PodSpec.SchedulerName == "" {
			bgd.Spec.PodSpec.SchedulerName = "default-scheduler"
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to fill unfilled pod spec fields of BlueGreenDeployment object with default values: %v", err)
	}

	return b, nil
}
