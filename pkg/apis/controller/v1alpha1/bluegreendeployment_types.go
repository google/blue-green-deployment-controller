package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BlueGreenDeploymentSpec defines the desired state of BlueGreenDeployment.
type BlueGreenDeploymentSpec struct {
	// Replicas is the number of desired replicas. It determines number of pods for both
	// Blue and Green ReplicaSets. It defaults to 1 if unspecified.
	// +optional
	Replicas int32 `json:"replicas,omitempty"`

	// PodSpec defines the specification of the desired behavior of the pods maintained
	// by the Blue and Green ReplicaSets, such as what image to run.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
	// +optional
	PodSpec v1.PodSpec `json:"podSpec,omitempty"`
}

// BlueGreenDeploymentStatus defines the observed state of BlueGreenDeployment
type BlueGreenDeploymentStatus struct {
	// ActiveReplicaSetColor records color of current active ReplicaSet, which is used to
	// change label selectors of service. It controls which ReplicaSet the service will
	// send traffic to. It defaults to "blue" color at the beginning.
	// +optional
	ActiveReplicaSetColor string `json:"activeReplicaSetColor,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BlueGreenDeployment
// +k8s:openapi-gen=true
// +resource:path=bluegreendeployments
type BlueGreenDeployment struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the specification of the desired behavior of the BlueGreenDeployment object.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
	// +optional
	Spec BlueGreenDeploymentSpec `json:"spec,omitempty"`

	// Status is the most recently observed status of the BlueGreenDeployment object.
	// This data may be out of date by some window of time.
	// Populated by the system.
	// Read-only.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
	// +optional
	Status BlueGreenDeploymentStatus `json:"status,omitempty"`
}
