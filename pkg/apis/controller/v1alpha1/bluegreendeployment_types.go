package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BlueGreenDeploymentSpec defines the desired state of BlueGreenDeployment
type BlueGreenDeploymentSpec struct {
	Replicas int32      `json:"replicas,omitempty"`
	PodSpec  v1.PodSpec `json:"podSpec,omitempty"`
}

// BlueGreenDeploymentStatus defines the observed state of BlueGreenDeployment
type BlueGreenDeploymentStatus struct {
	// record current active ReplicaSet, which is used to change label selectors of service
	ActiveReplicaSetColor string `json:"activeReplicaSetColor,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BlueGreenDeployment
// +k8s:openapi-gen=true
// +resource:path=bluegreendeployments
type BlueGreenDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BlueGreenDeploymentSpec   `json:"spec,omitempty"`
	Status BlueGreenDeploymentStatus `json:"status,omitempty"`
}
