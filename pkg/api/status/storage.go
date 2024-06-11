package status

import "encoding/json"

type StorageClassStatus struct {
	// currently no status
}

func (s *StorageClassStatus) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (s *StorageClassStatus) JsonMarshal() ([]byte, error) {
	return json.Marshal(s)
}

type PersistentVolumeClaimPhase string

const (
	PersistentVolumeClaimPending PersistentVolumeClaimPhase = "Pending"
	PersistentVolumeClaimBound   PersistentVolumeClaimPhase = "Bound"
	PersistentVolumeClaimLost    PersistentVolumeClaimPhase = "Lost"
)

type PersistentVolumeClaimStatus struct {
	Phase PersistentVolumeClaimPhase `json:"phase,omitempty"`
}

func (p *PersistentVolumeClaimStatus) JsonUnmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, p)
}

func (p *PersistentVolumeClaimStatus) JsonMarshal() ([]byte, error) {
	return json.Marshal(p)
}

type PersistentVolumePhase string

const (
	PersistentVolumeAvailable PersistentVolumePhase = "Available"
	PersistentVolumeBound     PersistentVolumePhase = "Bound"
	PersistentVolumeReleased  PersistentVolumePhase = "Released"
	PersistentVolumeFailed    PersistentVolumePhase = "Failed"
)

type PersistentVolumeStatus struct {
	Phase PersistentVolumePhase `json:"phase,omitempty"`
}

func (p *PersistentVolumeStatus) JsonUnmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, p)
}

func (p *PersistentVolumeStatus) JsonMarshal() ([]byte, error) {
	return json.Marshal(p)
}
