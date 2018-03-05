/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	demov1 "k8s.io/bgd-controller/pkg/apis/demo/v1"
	"k8s.io/client-go/kubernetes"
	typedv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/util/retry"
)

func newReplicaSet(name, color string, obj *demov1.BlueGreenDeployment, isInactive bool) *extensionsv1beta1.ReplicaSet {
	replicas := int32(obj.Spec.Replicas)
	if isInactive {
		replicas = int32(0)
	}
	return &extensionsv1beta1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ReplicaSet",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: obj.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(obj, schema.GroupVersionKind{
					Group:   demov1.SchemeGroupVersion.Group,
					Version: demov1.SchemeGroupVersion.Version,
					Kind:    "BlueGreenDeployment",
				}),
			},
			Annotations: map[string]string{
				bgdPodTemplateSpecAnnotation: "",
			},
		},
		Spec: extensionsv1beta1.ReplicaSetSpec{
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
							Name:  obj.Spec.Template.Spec.Containers[0].Name,
							Image: obj.Spec.Template.Spec.Containers[0].Image,
						},
					},
				},
			},
		},
	}
}

func newService(namespace string) *corev1.Service {
	labels := map[string]string{"color": "blue"}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "core/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bgd-svc",
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
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

// waitAllPodsAvailable returns true if all pods are available, false otherwise
func waitAllPodsAvailable(rs *extensionsv1beta1.ReplicaSet, kubeclient *kubernetes.Clientset) bool {
	if err := wait.PollImmediate(pollInterval, pollTimeout, func() (bool, error) {
		newRS, err := kubeclient.ExtensionsV1beta1().ReplicaSets(rs.Namespace).Get(rs.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return newRS.Status.Replicas == *rs.Spec.Replicas && newRS.Status.AvailableReplicas == *rs.Spec.Replicas, nil
	}); err != nil {
		fmt.Printf("failed to wait for all pods of replicaset %q to be available: %v\n", rs.Name, err)
		return false
	}
	return true
}

func scaleRS(rs *extensionsv1beta1.ReplicaSet, replicas int32, rsClient typedv1beta1.ReplicaSetInterface) (*extensionsv1beta1.ReplicaSet, error) {
	return updateRS(rsClient, rs.Name, func(rs *extensionsv1beta1.ReplicaSet) {
		*rs.Spec.Replicas = replicas
	})
}

func updateRS(rsClient typedv1beta1.ReplicaSetInterface, rsName string, updateFunc func(*extensionsv1beta1.ReplicaSet)) (*extensionsv1beta1.ReplicaSet, error) {
	var rs *extensionsv1beta1.ReplicaSet
	if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		newRS, err := rsClient.Get(rsName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		updateFunc(newRS)
		rs, err = rsClient.Update(newRS)
		return err
	}); err != nil {
		return nil, fmt.Errorf("Failed to update rs %s: %v", rsName, err)
	}
	return rs, nil
}

func waitRSDeletionToComplete(rs *extensionsv1beta1.ReplicaSet, rsClient typedv1beta1.ReplicaSetInterface) error {
	if err := wait.PollImmediate(pollInterval, pollTimeout, func() (bool, error) {
		_, err := rsClient.Get(rs.Name, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		return false, nil
	}); err != nil {
		return fmt.Errorf("failed to wait for replicaset %q deletion to complete: %v", rs.Name, err)
	}
	return nil
}

func updateService(svcName, namespace string, kubeClient *kubernetes.Clientset, updateFunc func(*corev1.Service)) (*corev1.Service, error) {
	var svc *corev1.Service
	svcClient := kubeClient.CoreV1().Services(namespace)
	if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		newSvc, err := svcClient.Get(svcName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		updateFunc(newSvc)
		svc, err = svcClient.Update(newSvc)
		return err
	}); err != nil {
		return nil, err
	}
	return svc, nil
}
