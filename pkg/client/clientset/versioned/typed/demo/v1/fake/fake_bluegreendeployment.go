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

package fake

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	demo_v1 "k8s.io/bgd-controller/pkg/apis/demo/v1"
	testing "k8s.io/client-go/testing"
)

// FakeBlueGreenDeployments implements BlueGreenDeploymentInterface
type FakeBlueGreenDeployments struct {
	Fake *FakeDemoV1
	ns   string
}

var bluegreendeploymentsResource = schema.GroupVersionResource{Group: "demo.google.com", Version: "v1", Resource: "bluegreendeployments"}

var bluegreendeploymentsKind = schema.GroupVersionKind{Group: "demo.google.com", Version: "v1", Kind: "BlueGreenDeployment"}

// Get takes name of the blueGreenDeployment, and returns the corresponding blueGreenDeployment object, and an error if there is any.
func (c *FakeBlueGreenDeployments) Get(name string, options v1.GetOptions) (result *demo_v1.BlueGreenDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(bluegreendeploymentsResource, c.ns, name), &demo_v1.BlueGreenDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*demo_v1.BlueGreenDeployment), err
}

// List takes label and field selectors, and returns the list of BlueGreenDeployments that match those selectors.
func (c *FakeBlueGreenDeployments) List(opts v1.ListOptions) (result *demo_v1.BlueGreenDeploymentList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(bluegreendeploymentsResource, bluegreendeploymentsKind, c.ns, opts), &demo_v1.BlueGreenDeploymentList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &demo_v1.BlueGreenDeploymentList{}
	for _, item := range obj.(*demo_v1.BlueGreenDeploymentList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested blueGreenDeployments.
func (c *FakeBlueGreenDeployments) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(bluegreendeploymentsResource, c.ns, opts))

}

// Create takes the representation of a blueGreenDeployment and creates it.  Returns the server's representation of the blueGreenDeployment, and an error, if there is any.
func (c *FakeBlueGreenDeployments) Create(blueGreenDeployment *demo_v1.BlueGreenDeployment) (result *demo_v1.BlueGreenDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(bluegreendeploymentsResource, c.ns, blueGreenDeployment), &demo_v1.BlueGreenDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*demo_v1.BlueGreenDeployment), err
}

// Update takes the representation of a blueGreenDeployment and updates it. Returns the server's representation of the blueGreenDeployment, and an error, if there is any.
func (c *FakeBlueGreenDeployments) Update(blueGreenDeployment *demo_v1.BlueGreenDeployment) (result *demo_v1.BlueGreenDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(bluegreendeploymentsResource, c.ns, blueGreenDeployment), &demo_v1.BlueGreenDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*demo_v1.BlueGreenDeployment), err
}

// Delete takes name of the blueGreenDeployment and deletes it. Returns an error if one occurs.
func (c *FakeBlueGreenDeployments) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(bluegreendeploymentsResource, c.ns, name), &demo_v1.BlueGreenDeployment{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBlueGreenDeployments) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(bluegreendeploymentsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &demo_v1.BlueGreenDeploymentList{})
	return err
}

// Patch applies the patch and returns the patched blueGreenDeployment.
func (c *FakeBlueGreenDeployments) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *demo_v1.BlueGreenDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(bluegreendeploymentsResource, c.ns, name, data, subresources...), &demo_v1.BlueGreenDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*demo_v1.BlueGreenDeployment), err
}
