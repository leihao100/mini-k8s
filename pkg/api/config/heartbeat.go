package config

import (
	"MiniK8S/pkg/api/meta"
	apitypes "MiniK8S/pkg/api/types"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const HeartbeatSendInterval = 10 * time.Second
const HeartbeatCheckInterval = 5 * time.Second
const HeartbeatTimeoutInterval = 40 * time.Second

type Heartbeat struct {
	ApiVersion string          `yaml:"apiVersion" json:"apiVersion,omitempty"`
	Kind       string          `yaml:"kind" json:"kind,omitempty"`
	Metadata   meta.ObjectMeta `json:"metadata,omitempty"`
	Uid        uuid.UUID       `json:"uid,omitempty"`
}

type HeartbeatList struct {
	ApiVersion      string      `json:"apiVersion,omitempty"`
	Kind            string      `json:"kind,omitempty"`
	ResourceVersion string      `json:"resourceVersion,omitempty"`
	Continue        string      `json:"continue,omitempty"`
	Items           []Heartbeat `json:"items"`
}

func (p *Heartbeat) JsonMarshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Heartbeat) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &p)
}
func (p *Heartbeat) GetName() string {
	return p.Metadata.Name
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

func (d *HeartbeatList) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &d)
}

func (d *HeartbeatList) JsonMarshal() ([]byte, error) {
	return json.Marshal(d)
}
func (d *HeartbeatList) AppendItems(objects []string) error {
	for _, object := range objects {
		ApiObject := &Heartbeat{}
		err := ApiObject.JsonUnmarshal([]byte(object))
		if err != nil {
			return err
		}
		d.Items = append(d.Items, *ApiObject)
	}
	return nil
}
func (d *HeartbeatList) GetItems() []ApiObject {
	var items []ApiObject
	items = make([]ApiObject, 0)
	for _, item := range d.Items {
		items = append(items, &item)
	}
	return items
}
func (d *HeartbeatList) Info() {
	fmt.Printf("%-10s\t%-10s\t%10s\t%-20s\n", "NAME", "UID", "DESIRED", "CURRENT")
	//for _, item := range d.Items {
	//fmt.Printf("%-10s\t%-10s\t%-10d\t%-20d\n", item.Metadata.Name, item.Metadata.Uid, item.Spec.Replicas, item.Status.Replicas)
	//}
}
