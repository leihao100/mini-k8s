package config

import (
	"MiniK8S/pkg/api/meta"
	apitypes "MiniK8S/pkg/api/types"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

const HeartbeatSendInterval = 10 * time.Second
const HeartbeatCheckInterval = 5 * time.Second
const HeartbeatTimeoutInterval = 40 * time.Second

type Heartbeat struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   meta.ObjectMeta
	Uid        uuid.UUID
}

func (p *Heartbeat) JsonMarshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Heartbeat) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &p)
}

func (p *Heartbeat) SetUID(uid uuid.UUID) {
	p.Metadata.Uid = uid
}

func (p *Heartbeat) GetUID() uuid.UUID {
	return p.Metadata.Uid
}

func (p *Heartbeat) SetResourceVersion(version int64) {
	p.Metadata.ResourceVersion = strconv.FormatInt(version, 10)
}
func (p *Heartbeat) GetResourceVersion() int64 {
	res, err := strconv.ParseInt(p.Metadata.ResourceVersion, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	return res
}
func (p *Heartbeat) JsonUnmarshalStatus(data []byte) error {
	//return json.Unmarshal(data, &(p.Status))
	return json.Unmarshal(data, &p)
}

func (p *Heartbeat) JsonMarshalStatus() ([]byte, error) {
	//return json.Marshal(p.Status)
	return json.Marshal(&p)
}
func (p *Heartbeat) SetStatus(s ApiObjectStatus) bool {
	//status, ok := s.(*status.PodStatus)
	//if ok {
	//	p.Status = *status
	//}
	return true
}
func (p *Heartbeat) GetStatus() ApiObjectStatus {
	return NewApiObjectStatus(apitypes.HeartbeatObjectType)
}
func (p *Heartbeat) Info() {
	//fmt.Printf("%-10s\t%-10s\t%-10s\t%-20s\t%-20s\n", "NAME", "UID", "NODE", "STATUS", "IP")
	//fmt.Printf("%-10s\t%-10s\t%-10s\t%-20s\t%-20s\n", p.Metadata.Name, p.Metadata.Uid, p.Spec.NodeName, p.Status.Phase, p.Status.PodIP)
}
