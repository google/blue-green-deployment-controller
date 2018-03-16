package fake

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	v1alpha1 "k8s.io/bgd/pkg/apis/controller/v1alpha1"
	testing "k8s.io/client-go/testing"
)

// FakeBlueGreenDeployments implements BlueGreenDeploymentInterface
type FakeBlueGreenDeployments struct {
	Fake *FakeControllerV1alpha1
	ns   string
}

var bluegreendeploymentsResource = schema.GroupVersionResource{Group: "controller.google.com", Version: "v1alpha1", Resource: "bluegreendeployments"}

var bluegreendeploymentsKind = schema.GroupVersionKind{Group: "controller.google.com", Version: "v1alpha1", Kind: "BlueGreenDeployment"}

// Get takes name of the blueGreenDeployment, and returns the corresponding blueGreenDeployment object, and an error if there is any.
func (c *FakeBlueGreenDeployments) Get(name string, options v1.GetOptions) (result *v1alpha1.BlueGreenDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(bluegreendeploymentsResource, c.ns, name), &v1alpha1.BlueGreenDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BlueGreenDeployment), err
}

// List takes label and field selectors, and returns the list of BlueGreenDeployments that match those selectors.
func (c *FakeBlueGreenDeployments) List(opts v1.ListOptions) (result *v1alpha1.BlueGreenDeploymentList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(bluegreendeploymentsResource, bluegreendeploymentsKind, c.ns, opts), &v1alpha1.BlueGreenDeploymentList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.BlueGreenDeploymentList{}
	for _, item := range obj.(*v1alpha1.BlueGreenDeploymentList).Items {
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
func (c *FakeBlueGreenDeployments) Create(blueGreenDeployment *v1alpha1.BlueGreenDeployment) (result *v1alpha1.BlueGreenDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(bluegreendeploymentsResource, c.ns, blueGreenDeployment), &v1alpha1.BlueGreenDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BlueGreenDeployment), err
}

// Update takes the representation of a blueGreenDeployment and updates it. Returns the server's representation of the blueGreenDeployment, and an error, if there is any.
func (c *FakeBlueGreenDeployments) Update(blueGreenDeployment *v1alpha1.BlueGreenDeployment) (result *v1alpha1.BlueGreenDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(bluegreendeploymentsResource, c.ns, blueGreenDeployment), &v1alpha1.BlueGreenDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BlueGreenDeployment), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBlueGreenDeployments) UpdateStatus(blueGreenDeployment *v1alpha1.BlueGreenDeployment) (*v1alpha1.BlueGreenDeployment, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(bluegreendeploymentsResource, "status", c.ns, blueGreenDeployment), &v1alpha1.BlueGreenDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BlueGreenDeployment), err
}

// Delete takes name of the blueGreenDeployment and deletes it. Returns an error if one occurs.
func (c *FakeBlueGreenDeployments) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(bluegreendeploymentsResource, c.ns, name), &v1alpha1.BlueGreenDeployment{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBlueGreenDeployments) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(bluegreendeploymentsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.BlueGreenDeploymentList{})
	return err
}

// Patch applies the patch and returns the patched blueGreenDeployment.
func (c *FakeBlueGreenDeployments) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.BlueGreenDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(bluegreendeploymentsResource, c.ns, name, data, subresources...), &v1alpha1.BlueGreenDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BlueGreenDeployment), err
}
