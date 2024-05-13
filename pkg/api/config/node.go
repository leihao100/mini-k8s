package config

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"time"
)

type Node struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   meta.ObjectMeta
	Spec       NodeSpec
	Status     status.NodeStatus
}

/*
NodeSpec from api document
已移除的部分已删去
podCIDR
string	PodCIDR represents the pod IP range assigned to the node.
podCIDRs
string array
patch strategy: merge	podCIDRs represents the IP ranges assigned to the node for usage by Pods on that node. If this field is specified, the 0th entry must match the podCIDR field. It may contain at most 1 value for each of IPv4 and IPv6.
providerID
string	ID of the node assigned by the cloud provider in the format: <ProviderName>://<ProviderSpecificNodeID>
taints
Taint array	If specified, the node's taints.
unschedulable
boolean	Unschedulable controls node schedulability of new pods. By default, node is schedulable. More info: https://kubernetes.io/docs/concepts/nodes/node/#manual-node-administration
*/
type NodeSpec struct {
	PodCIDR       string   `yaml:"podCIDR"`
	PodCIDRs      []string `yaml:"podCIDRs"`
	ProviderID    string   `yaml:"providerID"`
	Unschedulable bool     `yaml:"unschedulable"`
	Taints        []Taint  `yaml:"taints"`
}

type Taint struct {
	Key       string    `yaml:"key"`
	Value     string    `yaml:"value"`
	Effect    string    `yaml:"effect"`
	TimeAdded time.Time `yaml:"timeAdded"`
}
