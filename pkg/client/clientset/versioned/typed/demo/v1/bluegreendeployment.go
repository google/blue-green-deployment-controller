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
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	v1 "k8s.io/bgd-controller/pkg/apis/demo/v1"
	scheme "k8s.io/bgd-controller/pkg/client/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

// BlueGreenDeploymentsGetter has a method to return a BlueGreenDeploymentInterface.
// A group's client should implement this interface.
type BlueGreenDeploymentsGetter interface {
	BlueGreenDeployments(namespace string) BlueGreenDeploymentInterface
}

// BlueGreenDeploymentInterface has methods to work with BlueGreenDeployment resources.
type BlueGreenDeploymentInterface interface {
	Create(*v1.BlueGreenDeployment) (*v1.BlueGreenDeployment, error)
	Update(*v1.BlueGreenDeployment) (*v1.BlueGreenDeployment, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.BlueGreenDeployment, error)
	List(opts meta_v1.ListOptions) (*v1.BlueGreenDeploymentList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.BlueGreenDeployment, err error)
	BlueGreenDeploymentExpansion
}

// blueGreenDeployments implements BlueGreenDeploymentInterface
type blueGreenDeployments struct {
	client rest.Interface
	ns     string
}

// newBlueGreenDeployments returns a BlueGreenDeployments
func newBlueGreenDeployments(c *DemoV1Client, namespace string) *blueGreenDeployments {
	return &blueGreenDeployments{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the blueGreenDeployment, and returns the corresponding blueGreenDeployment object, and an error if there is any.
func (c *blueGreenDeployments) Get(name string, options meta_v1.GetOptions) (result *v1.BlueGreenDeployment, err error) {
	result = &v1.BlueGreenDeployment{}
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
func (c *blueGreenDeployments) List(opts meta_v1.ListOptions) (result *v1.BlueGreenDeploymentList, err error) {
	result = &v1.BlueGreenDeploymentList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested blueGreenDeployments.
func (c *blueGreenDeployments) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a blueGreenDeployment and creates it.  Returns the server's representation of the blueGreenDeployment, and an error, if there is any.
func (c *blueGreenDeployments) Create(blueGreenDeployment *v1.BlueGreenDeployment) (result *v1.BlueGreenDeployment, err error) {
	result = &v1.BlueGreenDeployment{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		Body(blueGreenDeployment).
		Do().
		Into(result)
	return
}

// Update takes the representation of a blueGreenDeployment and updates it. Returns the server's representation of the blueGreenDeployment, and an error, if there is any.
func (c *blueGreenDeployments) Update(blueGreenDeployment *v1.BlueGreenDeployment) (result *v1.BlueGreenDeployment, err error) {
	result = &v1.BlueGreenDeployment{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		Name(blueGreenDeployment.Name).
		Body(blueGreenDeployment).
		Do().
		Into(result)
	return
}

// Delete takes name of the blueGreenDeployment and deletes it. Returns an error if one occurs.
func (c *blueGreenDeployments) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *blueGreenDeployments) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bluegreendeployments").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched blueGreenDeployment.
func (c *blueGreenDeployments) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.BlueGreenDeployment, err error) {
	result = &v1.BlueGreenDeployment{}
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
