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

package v1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BlueGreenDeployment is a specification for a BlueGreenDeployment resource
type BlueGreenDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              BlueGreenDeploymentSpec   `json:"spec"`
	Status            BlueGreenDeploymentStatus `json:"status,omitempty"`
}

// BlueGreenDeploymentSpec is the spec for a BlueGreenDeployment resource
type BlueGreenDeploymentSpec struct {
	Replicas int32              `json:"replicas"`
	Template v1.PodTemplateSpec `json:"template"`
}

// BlueGreenDeploymentStatus is the status for a BlueGreenDeployment resource
type BlueGreenDeploymentStatus struct {
	Phase string `json:"phase,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BlueGreenDeploymentList is a list of BlueGreenDeployment resources
type BlueGreenDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []BlueGreenDeployment `json:"items"`
}
