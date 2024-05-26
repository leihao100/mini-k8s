package config

type ApiObject interface {
}
type ApiObjectSpec interface {
}
type ApiObjectStatus interface {
}

type ErrorApiObject struct {
}

type ApiObjectList interface {
	JsonUnmarshal(data []byte) error
	JsonMarshal() ([]byte, error)
	AddItemFromStr(objectStr string) error
	AppendItemsFromStr(objectStrs []string) error
	GetItems() any
	GetIApiObjectArr() []ApiObject
	PrintBrief()
}

// ListOptions is the query options to a standard REST list call.
type ListOptions struct {
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`

	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`

	// A selector to restrict the list of returned objects by their labels.
	// Defaults to everything.
	// +optional
	LabelSelector string `json:"labelSelector,omitempty" protobuf:"bytes,1,opt,name=labelSelector"`
	// A selector to restrict the list of returned objects by their fields.
	// Defaults to everything.
	// +optional
	FieldSelector string `json:"fieldSelector,omitempty" protobuf:"bytes,2,opt,name=fieldSelector"`

	// +k8s:deprecated=includeUninitialized,protobuf=6

	// Watch for changes to the described resources and return them as a stream of
	// add, update, and remove notifications. Specify resourceVersion.
	// +optional
	Watch bool `json:"watch,omitempty" protobuf:"varint,3,opt,name=watch"`

	// resourceVersion sets a constraint on what resource versions a request may be served from.
	// See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for
	// details.
	//
	// Defaults to unset
	// +optional
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,4,opt,name=resourceVersion"`

	// Timeout for the list/watch call.
	// This limits the duration of the call, regardless of any activity or inactivity.
	// +optional
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty" protobuf:"varint,5,opt,name=timeoutSeconds"`
}
