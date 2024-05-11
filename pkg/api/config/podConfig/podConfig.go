package podConfig

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/spec"
	"MiniK8S/pkg/api/status"
)

type PodConfig struct {
	ApiVersion string
	Kind       string
	Meta       meta.ObjectMeta
	Spec       spec.PodSpec
	Status     status.PodStatus
}

/*
API文档中描述如下
Field	Description
apiVersion
string	APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
kind
string	Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
metadata
ObjectMeta	Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
spec
PodSpec	Specification of the desired behavior of the pod. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
status
PodStatus	Most recently observed status of the pod. This data may not be up to date. Populated by the system. Read-only. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
*/
