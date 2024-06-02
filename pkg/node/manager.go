package node

import (
	"MiniK8S/pkg/api/address"
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/utils/net"
	"github.com/docker/docker/testutil"
	"github.com/google/uuid"
)

type NodeManager struct {
	node   *config.Node
	ty     config.NodeType
	Client *apiClient.Client
}

func CreateWorkerNode(n string) *NodeManager {
	Cli := apiClient.NewRESTClient(types.NodeObjectType)
	nc := &NodeManager{
		Client: Cli,
		node:   nil,
		ty:     config.Worker,
	}
	nc.Init(n)
	return nc
}

func CreateMasterNode(n string) *NodeManager {
	Cli := apiClient.NewRESTClient(types.NodeObjectType)
	nc := &NodeManager{
		Client: Cli,
		node:   nil,
		ty:     config.Master,
	}
	nc.Init(n)
	return nc

}

func (nm *NodeManager) GetNode() *config.Node {
	return nm.node
}

func (nm *NodeManager) Init(n string) {
	var name string
	switch nm.ty {
	case config.Worker:
		name = "worker-" + n
	default:
		name = "Master"
	}
	ip, _ := net.GetLocalIP()
	node := &config.Node{
		ApiVersion: "",
		Kind:       "node",
		Metadata: meta.ObjectMeta{
			Name: name,
			Uid:  uuid.New(),
		},
		Spec: config.NodeSpec{},
		Status: status.NodeStatus{
			Phase: "Pending",
			Addresses: address.NodeAddress{
				Address: ip,
				Type:    "ipv4",
			},
		},
	}
	nm.node = node
	url := nm.Client.BuildURL(apiClient.Create)
	buf, err := node.JsonMarshal()
	if err != nil {
		panic(err)
	}
	resp := nm.Client.Put(url, buf)
	if resp == nil {
		//error
	}
	nm.node.Status.Phase = "Running"
	url = nm.Client.BuildURL(apiClient.Create)
	buf, err = node.JsonMarshal()
	if err != nil {
		panic(err)
	}
	resp = nm.Client.Put(url, buf)
	if resp == nil {
		//error
	}
}

func GenerateRandomString(length int) string {
	return testutil.GenerateRandomAlphaOnlyString(length)
}
