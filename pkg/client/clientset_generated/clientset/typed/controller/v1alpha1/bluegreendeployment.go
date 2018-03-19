package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	v1alpha1 "k8s.io/bgd/pkg/apis/controller/v1alpha1"
	scheme "k8s.io/bgd/pkg/client/clientset_generated/clientset/scheme"
	rest "k8s.io/client-go/rest"
)

// BlueGreenDeploymentsGetter has a method to return a BlueGreenDeploymentInterface.
// A group's client should implement this interface.
type BlueGreenDeploymentsGetter interface {
	BlueGreenDeployments(namespace string) BlueGreenDeploymentInterface
}

// BlueGreenDeploymentInterface has methods to work with BlueGreenDeployment resources.
type BlueGreenDeploymentInterface interface {
	Create(*v1alpha1.BlueGreenDeployment) (*v1alpha1.BlueGreenDeployment, error)
	Update(*v1alpha1.BlueGreenDeployment) (*v1alpha1.BlueGreenDeployment, error)
	UpdateStatus(*v1alpha1.BlueGreenDeployment) (*v1alpha1.BlueGreenDeployment, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.BlueGreenDeployment, error)
	List(opts v1.ListOptions) (*v1alpha1.BlueGreenDeploymentList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.BlueGreenDeployment, err error)
	BlueGreenDeploymentExpansion
}

// blueGreenDeployments implements BlueGreenDeploymentInterface
type blueGreenDeployments struct {
	client rest.Interface
	ns     string
}

// newBlueGreenDeployments returns a BlueGreenDeployments
func newBlueGreenDeployments(c *ControllerV1alpha1Client, namespace string) *blueGreenDeployments {
	return &blueGreenDeployments{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the blueGreenDeployment, and returns the corresponding blueGreenDeployment object, and an error if there is any.
func (c *blueGreenDeployments) Get(name string, options v1.GetOptions) (result *v1alpha1.BlueGreenDeployment, err error) {
	result = &v1alpha1.BlueGreenDeployment{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BlueGreenDeployments that match those selectors.
func (c *blueGreenDeployments) List(opts v1.ListOptions) (result *v1alpha1.BlueGreenDeploymentList, err error) {
	result = &v1alpha1.BlueGreenDeploymentList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested blueGreenDeployments.
func (c *blueGreenDeployments) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a blueGreenDeployment and creates it.  Returns the server's representation of the blueGreenDeployment, and an error, if there is any.
func (c *blueGreenDeployments) Create(blueGreenDeployment *v1alpha1.BlueGreenDeployment) (result *v1alpha1.BlueGreenDeployment, err error) {
	result = &v1alpha1.BlueGreenDeployment{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		Body(blueGreenDeployment).
		Do().
		Into(result)
	return
}

// Update takes the representation of a blueGreenDeployment and updates it. Returns the server's representation of the blueGreenDeployment, and an error, if there is any.
func (c *blueGreenDeployments) Update(blueGreenDeployment *v1alpha1.BlueGreenDeployment) (result *v1alpha1.BlueGreenDeployment, err error) {
	result = &v1alpha1.BlueGreenDeployment{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		Name(blueGreenDeployment.Name).
		Body(blueGreenDeployment).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *blueGreenDeployments) UpdateStatus(blueGreenDeployment *v1alpha1.BlueGreenDeployment) (result *v1alpha1.BlueGreenDeployment, err error) {
	result = &v1alpha1.BlueGreenDeployment{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		Name(blueGreenDeployment.Name).
		SubResource("status").
		Body(blueGreenDeployment).
		Do().
		Into(result)
	return
}

// Delete takes name of the blueGreenDeployment and deletes it. Returns an error if one occurs.
func (c *blueGreenDeployments) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *blueGreenDeployments) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched blueGreenDeployment.
func (c *blueGreenDeployments) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.BlueGreenDeployment, err error) {
	result = &v1alpha1.BlueGreenDeployment{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("bluegreendeployments").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
