package config

import (
	"MiniK8S/pkg/api/meta"
	"encoding/json"
	"github.com/google/uuid"
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

func (hb *Heartbeat) JsonMarshal() ([]byte, error) {
	return json.Marshal(hb)
}

func (hb *Heartbeat) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &hb)
}

func (hb *Heartbeat) SetUID(uid uuid.UUID) {
	hb.Uid = uid
}

func (hb *Heartbeat) GetUID() uuid.UUID {
	return hb.Uid
}
