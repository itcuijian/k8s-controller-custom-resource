package v1

import (
	"github.com/itcuijim/k8s-controller-custom-resource/pkg/apis/samplecrd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemaGroupVersion = schema.GroupVersion{
	Group:   samplecrd.GroupName,
	Version: samplecrd.Version,
}

func Resource(resource string) schema.GroupResource {
	return SchemaGroupVersion.WithResource(resource).GroupResource()
}

func Kind(kind string) schema.GroupKind {
	return SchemaGroupVersion.WithKind(kind).GroupKind()
}

func addKnownTypes(schema *runtime.Scheme) error {
	schema.AddKnownTypes(
		SchemaGroupVersion,
		&Network{},
		&NetworkList{},
	)

	metav1.AddToGroupVersion(schema, SchemaGroupVersion)
	return nil
}
