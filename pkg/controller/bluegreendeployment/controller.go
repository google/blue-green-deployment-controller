package bluegreendeployment

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/kubernetes-sigs/kubebuilder/pkg/controller"
	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/types"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	controllerv1alpha1 "k8s.io/bgd/pkg/apis/controller/v1alpha1"
	controllerv1alpha1client "k8s.io/bgd/pkg/client/clientset_generated/clientset/typed/controller/v1alpha1"
	"k8s.io/bgd/pkg/inject/args"
	bgdutil "k8s.io/bgd/pkg/util"
	"k8s.io/client-go/kubernetes"
)

const (
	BlueColor               = "blue"
	GreenColor              = "green"
	BlueGreenDeploymentKind = "BlueGreenDeployment"
	BlueRSName              = BlueColor + "-rs"
	GreenRSName             = GreenColor + "-rs"
	ServiceName             = "bgd-svc"

	LabelColorKey = "color"

	PollInterval = 100 * time.Millisecond
	PollTimeout  = 30 * time.Second
)

// Reconcile reconciles BlueGreenDeployment object to desired state.
func (c *BlueGreenDeploymentController) Reconcile(k types.ReconcileKey) error {
	b, err := c.blueGreenDeploymentClient.BlueGreenDeployments(k.Namespace).Get(k.Name, metav1.GetOptions{})
	if err != nil || b == nil {
		return err
	}

	if b, err = defaultFields(c.blueGreenDeploymentClient, b); err != nil {
		return err
	}

	if err := c.sync(b); err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

// +controller:group=controller,version=v1alpha1,kind=BlueGreenDeployment,resource=bluegreendeployments
// +rbac:groups=controller,resources=bluegreendeployments,verbs=get;list;watch;create;update;patch;delete
type BlueGreenDeploymentController struct {
	// blueGreenDeploymentLister *controllerv1alpha1lister.BlueGreenDeploymentLister
	blueGreenDeploymentClient controllerv1alpha1client.ControllerV1alpha1Interface
	k8sClientSet              *kubernetes.Clientset
}

// ProvideController provides a controller that will be run at startup. Kubebuilder will use code generation
// to automatically register this controller in the inject package.
func ProvideController(arguments args.InjectArgs) (*controller.GenericController, error) {
	bc := &BlueGreenDeploymentController{
		blueGreenDeploymentClient: arguments.Clientset.ControllerV1alpha1(),
		k8sClientSet:              arguments.KubernetesClientSet,
	}

	// Create a new controller that will call BlueGreenDeploymentController.Reconcile on changes to BlueGreenDeployments
	gc := &controller.GenericController{
		Name:             BlueGreenDeploymentKind + "Controller",
		Reconcile:        bc.Reconcile,
		InformerRegistry: arguments.ControllerManager,
	}
	if err := gc.Watch(&controllerv1alpha1.BlueGreenDeployment{}); err != nil {
		return gc, err
	}
	var gvk = metav1.GroupVersionKind{Group: "apps", Version: "v1", Kind: "ReplicaSet"}
	if err := gc.WatchAndMapToController(&appsv1.ReplicaSet{}, gvk); err != nil {
		return gc, err
	}

	return gc, nil
}

// sync lies bulk of reconciliation logic:
// 1. create blue ReplicaSet if it is not found
// 2. create green ReplicaSet if it is not found
// 3. create service if it is not found
// 4. reconcile blue and green ReplicaSets based on BlueGreenDeployment object's pod spec
func (c *BlueGreenDeploymentController) sync(b *controllerv1alpha1.BlueGreenDeployment) error {
	var err error

	// Initialize .status.activeReplicaSetColor field to blue color if it is not set yet
	if b.Status.ActiveReplicaSetColor == "" {
		// Update .status.activeReplicaSetColor field to blue color at the beginning
		if b, err = updateBlueGreenDeployment(c.blueGreenDeploymentClient, b.Name, b.Namespace, func(bgd *controllerv1alpha1.BlueGreenDeployment) {
			bgd.Status.ActiveReplicaSetColor = BlueColor
		}); err != nil {
			return fmt.Errorf("failed to update .status.activeReplicaSetColor field of BlueGreenDeployment: %v", err)
		}
	}

	blueRS, greenRS, err := createReplicaSetsIfNotFound(b, c.k8sClientSet)
	if err != nil {
		return err
	}

	service, err := createServiceIfNotFound(b, c.k8sClientSet)
	if err != nil {
		return err
	}

	err = c.reconcileActiveReplicaSet(b, blueRS, greenRS, service)
	if err != nil {
		return err
	}

	return nil
}

// reconcileActiveReplicaSet checks active and inactive ReplicaSets' pod specs to
// see if any of the pod specs matches BlueGreenDeployment object's pod spec.
// If the active ReplicaSet has the matching spec, the controller does nothing.
// Else if the inactive ReplicaSet has the matching spec, the controller:
// 1. scales up the inactive ReplicaSet
// 2. modifies label selectors of the service to point to the newly active ReplicaSet
// 3. scales down previously active ReplicaSet to make it inactive
// Else (i.e., none of the Replicasets has the matching spec), the controller:
// 1. deletes the inactive ReplicaSet to make room for a new ReplicaSet
// 2. creates and scales up the new ReplicaSet
// 3. modifies label selectors of the service to point to the new ReplicaSet
// 4. scales down previously active ReplicaSet to make it inactive
// The function updates .status.activeReplicaSetColor field and the service's label
// selectors with color of the active ReplicaSet.
func (c *BlueGreenDeploymentController) reconcileActiveReplicaSet(b *controllerv1alpha1.BlueGreenDeployment, blueRS, greenRS *appsv1.ReplicaSet, service *corev1.Service) error {
	var err error

	// If a ReplicaSet has the same color as indicated by .status.activeReplicaSetColor field
	// of BlueGreenDeployment object, it is the active ReplicaSet.
	activeColor := b.Status.ActiveReplicaSetColor
	activeRS := blueRS
	inactiveRS := greenRS
	if activeColor == GreenColor {
		activeRS = greenRS
		inactiveRS = blueRS
	}

	// Hash pod spec of active ReplicaSet, inactive ReplicaSet, and BlueGreenDeployment object for comparison
	hashedActiveRSPodSpec := bgdutil.ComputeHash(&activeRS.Spec.Template.Spec)
	hashedInactiveRSPodSpec := bgdutil.ComputeHash(&inactiveRS.Spec.Template.Spec)
	hashedBGDPodSpec := bgdutil.ComputeHash(&b.Spec.PodSpec)

	// Case: active ReplicaSet has the matching spec
	if reflect.DeepEqual(hashedActiveRSPodSpec, hashedBGDPodSpec) {
		return nil
	} else {
		// Update .status.activeReplicaSetColor field to inactive ReplicaSet's color as it will become active soon
		if b, err = updateBlueGreenDeployment(c.blueGreenDeploymentClient, b.Name, b.Namespace, func(bgd *controllerv1alpha1.BlueGreenDeployment) {
			bgd.Status.ActiveReplicaSetColor = inactiveRS.Spec.Selector.MatchLabels[LabelColorKey]
		}); err != nil {
			return fmt.Errorf("failed to update .status.activeReplicaSetColor field of BlueGreenDeployment: %v", err)
		}

		// Case: inactive ReplicaSet has the matching spec
		if reflect.DeepEqual(hashedInactiveRSPodSpec, hashedBGDPodSpec) {
			err = scaleReplicaSet(c.k8sClientSet, b.Spec.Replicas, inactiveRS)
			if err != nil {
				return fmt.Errorf("during scaling up of inactive ReplicaSet %q, %v", inactiveRS.Name, err)
			}
		} else { // Case: no Replicaset has the matching spec
			_, err = replaceInactiveReplicaSet(c.k8sClientSet, b, inactiveRS)
			if err != nil {
				return fmt.Errorf("during replacement of inactive ReplicaSet %q, %v", inactiveRS.Name, err)
			}
		}

		// Scale down previously active ReplicaSet
		err = scaleReplicaSet(c.k8sClientSet, 0, activeRS)
		if err != nil {
			return fmt.Errorf("during scaling down of active ReplicaSet %q, %v", activeRS.Name, err)
		}

		// Point the service to current active ReplicaSet by updating its "color" label selector to match active ReplicaSet's color
		if service, err = updateService(c.k8sClientSet, service.Namespace, func(service *corev1.Service) {
			service.Spec.Selector[LabelColorKey] = b.Status.ActiveReplicaSetColor
		}); err != nil {
			return fmt.Errorf("failed to update %q label selector of service %q to match active ReplicaSet's color: %v", LabelColorKey, service.Name, err)
		}
	}

	return nil
}
