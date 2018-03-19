package v1alpha1

import (
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	v1alpha1 "k8s.io/bgd/pkg/apis/controller/v1alpha1"
	"k8s.io/bgd/pkg/client/clientset_generated/clientset/scheme"
	rest "k8s.io/client-go/rest"
)

type ControllerV1alpha1Interface interface {
	RESTClient() rest.Interface
	BlueGreenDeploymentsGetter
}

// ControllerV1alpha1Client is used to interact with features provided by the controller.google.com group.
type ControllerV1alpha1Client struct {
	restClient rest.Interface
}

func (c *ControllerV1alpha1Client) BlueGreenDeployments(namespace string) BlueGreenDeploymentInterface {
	return newBlueGreenDeployments(c, namespace)
}

// NewForConfig creates a new ControllerV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*ControllerV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &ControllerV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new ControllerV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *ControllerV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new ControllerV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *ControllerV1alpha1Client {
	return &ControllerV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *ControllerV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
