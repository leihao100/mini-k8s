package config

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type NodeType string

const (
	Master NodeType = "master"
	Worker NodeType = "worker"
)

type Node struct {
	ApiVersion string            `yaml:"apiVersion" json:"apiVersion,omitempty"`
	Kind       string            `yaml:"kind" json:"kind,omitempty"`
	Metadata   meta.ObjectMeta   `json:"metadata,omitempty"`
	Spec       NodeSpec          `json:"spec,omitempty"`
	Status     status.NodeStatus `json:"status,omitempty"`
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
	PodCIDR       string   `yaml:"podCIDR" json:"podCIDR,omitempty"`
	PodCIDRs      []string `yaml:"podCIDRs" json:"podCIDRs,omitempty"`
	ProviderID    string   `yaml:"providerID" json:"providerID,omitempty"`
	Unschedulable bool     `yaml:"unschedulable" json:"unschedulable,omitempty"`
	Taints        []Taint  `yaml:"taints" json:"taints,omitempty"`
}

type Taint struct {
	Key       string    `yaml:"key" json:"key,omitempty"`
	Value     string    `yaml:"value" json:"value,omitempty"`
	Effect    string    `yaml:"effect" json:"effect,omitempty"`
	TimeAdded time.Time `yaml:"timeAdded" json:"timeAdded,omitempty"`
}

type NodeList struct {
	ApiVersion      string `json:"apiVersion,omitempty" `
	Kind            string `json:"kind,omitempty"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
	Continue        string `json:"continue,omitempty"`
	Items           []Node `json:"items"`
}

func (n *Node) JsonMarshal() ([]byte, error) {
	return json.Marshal(n)
}

func (n *Node) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &n)
}

func (n *Node) SetUID(uid uuid.UUID) {
	n.Metadata.Uid = uid
}

func (n *Node) GetUID() uuid.UUID {
	return n.Metadata.Uid
}

func (n *Node) SetResourceVersion(version int64) {
	n.Metadata.ResourceVersion = strconv.FormatInt(version, 10)
}
func (n *Node) GetResourceVersion() int64 {
	res, err := strconv.ParseInt(n.Metadata.ResourceVersion, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	return res
}
func (n *Node) JsonUnmarshalStatus(data []byte) error {
	return json.Unmarshal(data, &(n.Status))
}

func (n *Node) JsonMarshalStatus() ([]byte, error) {
	return json.Marshal(n.Status)
}
func (n *Node) SetStatus(ss ApiObjectStatus) bool {
	status, ok := ss.(*status.NodeStatus)
	if ok {
		n.Status = *status
	}
	return ok
}
func (n *Node) GetStatus() ApiObjectStatus {
	return &n.Status
}

func (n *Node) Info() {
	fmt.Printf("%-10s\t%-10s\t%-10s\t%-20s\n", "NAME", "UID", "STATUS", "IP")
	fmt.Printf("%-10s\t%-10s\t%-10s\t%-20s\n", n.Metadata.Name, n.Metadata.Uid, n.Status.Phase, n.Status.Addresses.Address)
}
func (n *NodeList) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &n)
}

func (n *NodeList) JsonMarshal() ([]byte, error) {
	return json.Marshal(n)
}
func (n *NodeList) AppendItems(objects []string) error {
	for _, object := range objects {
		ApiObject := &Node{}
		err := ApiObject.JsonUnmarshal([]byte(object))
		if err != nil {
			return err
		}
		n.Items = append(n.Items, *ApiObject)
	}
	return nil
}
func (n *NodeList) GetItems() []ApiObject {
	var items []ApiObject
	items = make([]ApiObject, 0)
	for _, item := range n.Items {
		items = append(items, &item)
	}
	return items
}
func (n *NodeList) Info() {
	fmt.Printf("%-10s\t%-10s\t%10s\t%-20s\n", "NAME", "UID", "STATUS", "IP")
	for _, item := range n.Items {
		fmt.Printf("%-10s\t%-10s\t%-10s\t%-20s\t\n", item.Metadata.Name, item.Metadata.Uid, item.Status.Phase, item.Status.Addresses.Address)
	}
}
