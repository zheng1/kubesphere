package main

import (
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Resource struct {
	Kind      string
	Namespace string
	Name      string
}

func GetResourceFromString(s string) Resource {
	rs := strings.Split(s, "/")
	return Resource{
		Kind:      rs[0],
		Namespace: rs[1],
		Name:      rs[2],
	}
}

func getNameFromObject(object runtime.RawExtension) string {
	return jsoniter.Get(object.Raw, "metadata", "name").ToString()
}

func GetResoureceFromAdmissionRequest(r *v1beta1.AdmissionRequest) Resource {
	resource := Resource{
		Kind:      r.Kind.Kind,
		Name:      getNameFromObject(r.Object),
		Namespace: r.Namespace,
	}
	return resource
}

func (resource Resource) GetId() string {
	return fmt.Sprintf("%s/%s/%s", resource.Kind, resource.Namespace, resource.Name)
}
