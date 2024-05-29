package status

import (
	"MiniK8S/pkg/api/address"
	"encoding/json"
)

type NodeStatus struct {
	Addresses address.NodeAddress `json:"addresses,omitempty"`
	//Allocatable object
	DaemonEndpoints int64  `json:"daemonEndpoints,omitempty"` //其对应了kubelet所监听的端口
	Phase           string `json:"phase,omitempty"`
}

/*
Field	Description
addresses
NodeAddress array
patch strategy: merge
patch merge key: type	List of addresses reachable to the node. Queried from cloud provider, if available. More info: https://kubernetes.io/docs/concepts/nodes/node/#addresses Note: This field is declared as mergeable, but the merge key is not sufficiently unique, which can cause data corruption when it is merged. Callers should instead use a full-replacement patch. See http://pr.k8s.io/79391 for an example.
allocatable
object	Allocatable represents the resources of a node that are available for scheduling. Defaults to Capacity.
capacity
object	Capacity represents the total resources of a node. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#capacity
conditions
NodeCondition array
patch strategy: merge
patch merge key: type	Conditions is an array of current observed node conditions. More info: https://kubernetes.io/docs/concepts/nodes/node/#condition
config
NodeConfigStatus	Status of the config assigned to the node via the dynamic Kubelet config feature.
daemonEndpoints
NodeDaemonEndpoints	Endpoints of daemons running on the Node.
images
ContainerImage array	List of container images on this node
nodeInfo
NodeSystemInfo	Set of ids/uuids to uniquely identify the node. More info: https://kubernetes.io/docs/concepts/nodes/node/#info
phase
string	NodePhase is the recently observed lifecycle phase of the node. More info: https://kubernetes.io/docs/concepts/nodes/node/#phase The field is never populated, and now is deprecated.
volumesAttached
AttachedVolume array	List of volumes that are attached to the node.
volumesInUse
string array	List of attachable volumes in use (mounted) by the node.
*/
func (n *NodeStatus) JsonMarshal() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NodeStatus) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &n)
}
