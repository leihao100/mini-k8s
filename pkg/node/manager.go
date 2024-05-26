package node

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"encoding/base64"
	"github.com/google/uuid"
	"math/rand"
)

type NodeManager struct {
	node   *config.Node
	ty     config.NodeType
	Client *apiClient.Client
}

func CreateWorkerNode() *NodeManager {
	Cli := apiClient.NewRESTClient(types.NodeObjectType)
	nc := &NodeManager{
		Client: Cli,
		node:   nil,
		ty:     config.Worker,
	}
	nc.Init()
	return nc
}

func CreateMasterNode() *NodeManager {
	Cli := apiClient.NewRESTClient(types.NodeObjectType)
	nc := &NodeManager{
		Client: Cli,
		node:   nil,
		ty:     config.Master,
	}
	nc.Init()
	return nc

}

func (nm *NodeManager) GetNode() *config.Node {
	return nm.node
}

func (nm *NodeManager) Init() {
	var name string
	switch nm.ty {
	case config.Worker:
		name = "worker-" + GenerateRandomString(5)
	default:
		name = "Master"
	}
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
		},
	}
	nm.node = node
	url := nm.Client.BuildURL(apiClient.Create)
	buf, err := node.JsonMarshal()
	if err != nil {
		panic(err)
	}
	resp := nm.Client.Post(url, buf)
	if resp == nil {
		//error
	}
	nm.node.Status.Phase = "Running"
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
