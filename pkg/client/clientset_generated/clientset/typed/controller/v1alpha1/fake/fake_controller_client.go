package fake

import (
	v1alpha1 "k8s.io/bgd/pkg/client/clientset_generated/clientset/typed/controller/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeControllerV1alpha1 struct {
	*testing.Fake
}

func (c *FakeControllerV1alpha1) BlueGreenDeployments(namespace string) v1alpha1.BlueGreenDeploymentInterface {
	return &FakeBlueGreenDeployments{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeControllerV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
