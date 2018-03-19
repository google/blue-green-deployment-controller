


// Api versions allow the api contract for a resource to be changed while keeping
// backward compatibility by support multiple concurrent versions
// of the same resource

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=k8s.io/bgd/pkg/apis/controller
// +k8s:defaulter-gen=TypeMeta
// +groupName=controller.google.com
package v1alpha1 // import "k8s.io/bgd/pkg/apis/controller/v1alpha1"
