// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "github.com/itcuijim/k8s-controller-custom-resource/pkg/apis/samplecrd/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeNetworks implements NetworkInterface
type FakeNetworks struct {
	Fake *FakeSamplecrdV1
	ns   string
}

var networksResource = v1.SchemeGroupVersion.WithResource("networks")

var networksKind = v1.SchemeGroupVersion.WithKind("Network")

// Get takes name of the network, and returns the corresponding network object, and an error if there is any.
func (c *FakeNetworks) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Network, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(networksResource, c.ns, name), &v1.Network{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Network), err
}

// List takes label and field selectors, and returns the list of Networks that match those selectors.
func (c *FakeNetworks) List(ctx context.Context, opts metav1.ListOptions) (result *v1.NetworkList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(networksResource, networksKind, c.ns, opts), &v1.NetworkList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.NetworkList{ListMeta: obj.(*v1.NetworkList).ListMeta}
	for _, item := range obj.(*v1.NetworkList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested networks.
func (c *FakeNetworks) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(networksResource, c.ns, opts))

}

// Create takes the representation of a network and creates it.  Returns the server's representation of the network, and an error, if there is any.
func (c *FakeNetworks) Create(ctx context.Context, network *v1.Network, opts metav1.CreateOptions) (result *v1.Network, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(networksResource, c.ns, network), &v1.Network{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Network), err
}

// Update takes the representation of a network and updates it. Returns the server's representation of the network, and an error, if there is any.
func (c *FakeNetworks) Update(ctx context.Context, network *v1.Network, opts metav1.UpdateOptions) (result *v1.Network, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(networksResource, c.ns, network), &v1.Network{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Network), err
}

// Delete takes name of the network and deletes it. Returns an error if one occurs.
func (c *FakeNetworks) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(networksResource, c.ns, name, opts), &v1.Network{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNetworks) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(networksResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1.NetworkList{})
	return err
}

// Patch applies the patch and returns the patched network.
func (c *FakeNetworks) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Network, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(networksResource, c.ns, name, pt, data, subresources...), &v1.Network{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Network), err
}
