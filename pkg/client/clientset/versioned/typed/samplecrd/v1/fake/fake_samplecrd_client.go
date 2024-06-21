// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "github.com/itcuijim/k8s-controller-custom-resource/pkg/client/clientset/versioned/typed/samplecrd/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeSamplecrdV1 struct {
	*testing.Fake
}

func (c *FakeSamplecrdV1) Networks(namespace string) v1.NetworkInterface {
	return &FakeNetworks{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeSamplecrdV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
