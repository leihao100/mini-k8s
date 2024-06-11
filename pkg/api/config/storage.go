package config

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"strconv"
)

type StorageClass struct {
	ApiVersion string                    `json:"apiVersion,omitempty"`
	Kind       string                    `json:"kind,omitempty"`
	Metadata   meta.ObjectMeta           `json:"metadata,omitempty"`
	Spec       StorageClassSpec          `json:"spec,omitempty"`
	Status     status.StorageClassStatus `json:"status,omitempty"`
}

type StorageClassSpec struct {
	// Provisioner currently only support cloud-os-2024-group-10.sjtu.edu.cn/nfs
	Provisioner string        `json:"provisioner,omitempty"`
	Parameters  NFSParameters `json:"parameters,omitempty"`
	// ReclaimPolicy currently only support Retain
	ReclaimPolicy string `json:"reclaimPolicy,omitempty"`
}

const (
	StorageClassNFSProvisioner      = "cloud-os-2024-group-10.sjtu.edu.cn/nfs"
	StorageClassRetainReclaimPolicy = "Retain"
)

type NFSParameters struct {
	Server string `json:"server,omitempty"`
	Path   string `json:"path,omitempty"`
}

func (s *StorageClass) JsonUnmarshal(bytes []byte) error {
	err := json.Unmarshal(bytes, s)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageClass) JsonMarshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *StorageClass) SetUID(uuid uuid.UUID) {
	s.Metadata.Uid = uuid
}

func (s *StorageClass) GetUID() uuid.UUID {
	return s.Metadata.Uid
}

func (s *StorageClass) GetName() string {
	return s.Metadata.Name
}

func (s *StorageClass) SetResourceVersion(version int64) {
	s.Metadata.ResourceVersion = strconv.FormatInt(version, 10)
}

func (s *StorageClass) GetResourceVersion() int64 {
	res, err := strconv.ParseInt(s.Metadata.ResourceVersion, 10, 64)
	if err != nil {
		log.Println("[GetResourceVersion] Error:", err)
		return 0
	}
	return res
}

func (s *StorageClass) JsonUnmarshalStatus(data []byte) error {
	return json.Unmarshal(data, &(s.Status))
}

func (s *StorageClass) JsonMarshalStatus() ([]byte, error) {
	return json.Marshal(s.Status)
}

func (s *StorageClass) SetStatus(st ApiObjectStatus) bool {
	sta, ok := st.(*status.StorageClassStatus)
	if ok {
		s.Status = *sta
	}
	return ok
}

func (s *StorageClass) GetStatus() ApiObjectStatus {
	return &s.Status
}

func (s *StorageClass) Info() {
	fmt.Printf("%-10s\t%-10s\t%10s\t%-20s\n", "NAME", "UID", "SERVER", "PATH")
	fmt.Printf("%-10s\t%-10s\t%10s\t%-20s\n", s.Metadata.Name, s.Metadata.Uid, s.Spec.Parameters.Server, s.Spec.Parameters.Path)
}

type StorageClassList struct {
	ApiVersion string         `json:"apiVersion,omitempty"`
	Kind       string         `json:"kind,omitempty"`
	Items      []StorageClass `json:"items"`
}

func (s *StorageClassList) JsonUnmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, s)
}

func (s *StorageClassList) JsonMarshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *StorageClassList) AppendItems(objects []string) error {
	for _, object := range objects {
		var sc StorageClass
		err := sc.JsonUnmarshal([]byte(object))
		if err != nil {
			return err
		}
		s.Items = append(s.Items, sc)
	}
	return nil
}

func (s *StorageClassList) GetItems() []ApiObject {
	var items []ApiObject
	items = make([]ApiObject, 0)
	for _, item := range s.Items {
		items = append(items, &item)
	}
	return items
}

func (s *StorageClassList) Info() {
	fmt.Printf("%-10s\t%-40s\t%20s\t%-10s\n", "NAME", "UID", "SERVER", "PATH")
	for _, item := range s.Items {
		fmt.Printf("%-10s\t%-40s\t%20s\t%-10s\n", item.Metadata.Name, item.Metadata.Uid, item.Spec.Parameters.Server, item.Spec.Parameters.Path)
	}
}

type PersistentVolume struct {
	ApiVersion string                        `json:"apiVersion,omitempty"`
	Kind       string                        `json:"kind,omitempty"`
	Metadata   meta.ObjectMeta               `json:"metadata,omitempty"`
	Spec       PersistentVolumeSpec          `json:"spec,omitempty"`
	Status     status.PersistentVolumeStatus `json:"status,omitempty"`
}

type PersistentVolumeSpec struct {
	ClaimRef         meta.OwnerReference `json:"claimRef,omitempty"`
	StorageClassName string              `json:"storageClassName,omitempty"`
	Nfs              NFSVolumeSource     `json:"nfs,omitempty"`
}

type NFSVolumeSource struct {
	Server string `json:"server,omitempty"`
	Path   string `json:"path,omitempty"`
}

func (p *PersistentVolume) JsonUnmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, p)
}

func (p *PersistentVolume) JsonMarshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PersistentVolume) SetUID(u uuid.UUID) {
	p.Metadata.Uid = u
}

func (p *PersistentVolume) GetUID() uuid.UUID {
	return p.Metadata.Uid
}

func (p *PersistentVolume) GetName() string {
	return p.Metadata.Name
}

func (p *PersistentVolume) SetResourceVersion(i int64) {
	p.Metadata.ResourceVersion = strconv.FormatInt(i, 10)
}

func (p *PersistentVolume) GetResourceVersion() int64 {
	res, err := strconv.ParseInt(p.Metadata.ResourceVersion, 10, 64)
	if err != nil {
		fmt.Println("[GetResourceVersion] Error:", err)
		return 0
	}
	return res
}

func (p *PersistentVolume) JsonUnmarshalStatus(bytes []byte) error {
	return json.Unmarshal(bytes, &(p.Status))
}

func (p *PersistentVolume) JsonMarshalStatus() ([]byte, error) {
	return json.Marshal(p.Status)
}

func (p *PersistentVolume) SetStatus(objectStatus ApiObjectStatus) bool {
	sta, ok := objectStatus.(*status.PersistentVolumeStatus)
	if ok {
		p.Status = *sta
	}
	return ok
}

func (p *PersistentVolume) GetStatus() ApiObjectStatus {
	return &p.Status
}

func (p *PersistentVolume) Info() {
	fmt.Printf("%-10s\t%-40s\t%-20s\t%-10s\t%-20s\t%-10s\t%-10s\n", "NAME", "UID", "OWNER", "SCLASS", "SERVER", "PATH", "PHASE")
	fmt.Printf("%-10s\t%-40s\t%-20s\t%-10s\t%-20s\t%-10s\t%-10s\n", p.Metadata.Name, p.Metadata.Uid, p.Spec.ClaimRef.Name, p.Spec.StorageClassName, p.Spec.Nfs.Server, p.Spec.Nfs.Path, p.Status.Phase)
}

type PersistentVolumeList struct {
	ApiVersion string             `json:"apiVersion,omitempty"`
	Kind       string             `json:"kind,omitempty"`
	Items      []PersistentVolume `json:"items"`
}

func (p *PersistentVolumeList) JsonUnmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, p)
}

func (p *PersistentVolumeList) JsonMarshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PersistentVolumeList) AppendItems(objects []string) error {
	for _, object := range objects {
		var pv PersistentVolume
		err := pv.JsonUnmarshal([]byte(object))
		if err != nil {
			return err
		}
		p.Items = append(p.Items, pv)
	}
	return nil
}

func (p *PersistentVolumeList) GetItems() []ApiObject {
	var items []ApiObject
	items = make([]ApiObject, 0)
	for _, item := range p.Items {
		items = append(items, &item)
	}
	return items
}

func (p *PersistentVolumeList) Info() {
	fmt.Printf("%-10s\t%-40s\t%-20s\t%-10s\t%-20s\t%-10s\t%-10s\n", "NAME", "UID", "OWNER", "SCLASS", "SERVER", "PATH", "PHASE")
	for _, item := range p.Items {
		fmt.Printf("%-10s\t%-40s\t%-20s\t%-10s\t%-20s\t%-10s\t%-10s\n", item.Metadata.Name, item.Metadata.Uid, item.Spec.ClaimRef.Name, item.Spec.StorageClassName, item.Spec.Nfs.Server, item.Spec.Nfs.Path, item.Status.Phase)
	}
}

type PersistentVolumeClaim struct {
	ApiVersion string                             `json:"apiVersion,omitempty"`
	Kind       string                             `json:"kind,omitempty"`
	Metadata   meta.ObjectMeta                    `json:"metadata,omitempty"`
	Spec       PersistentVolumeClaimSpec          `json:"spec,omitempty"`
	Status     status.PersistentVolumeClaimStatus `json:"status,omitempty"`
}

type PersistentVolumeClaimSpec struct {
	VolumeName       string `json:"volumeName,omitempty"`
	StorageClassName string `json:"storageClassName,omitempty"`
}

func (p *PersistentVolumeClaim) JsonUnmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, p)
}

func (p *PersistentVolumeClaim) JsonMarshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PersistentVolumeClaim) SetUID(u uuid.UUID) {
	p.Metadata.Uid = u
}

func (p *PersistentVolumeClaim) GetUID() uuid.UUID {
	return p.Metadata.Uid
}

func (p *PersistentVolumeClaim) GetName() string {
	return p.Metadata.Name
}

func (p *PersistentVolumeClaim) SetResourceVersion(i int64) {
	p.Metadata.ResourceVersion = strconv.FormatInt(i, 10)
}

func (p *PersistentVolumeClaim) GetResourceVersion() int64 {
	res, err := strconv.ParseInt(p.Metadata.ResourceVersion, 10, 64)
	if err != nil {
		fmt.Println("[GetResourceVersion] Error:", err)
		return 0
	}
	return res
}

func (p *PersistentVolumeClaim) JsonUnmarshalStatus(bytes []byte) error {
	return json.Unmarshal(bytes, &(p.Status))
}

func (p *PersistentVolumeClaim) JsonMarshalStatus() ([]byte, error) {
	return json.Marshal(p.Status)
}

func (p *PersistentVolumeClaim) SetStatus(objectStatus ApiObjectStatus) bool {
	sta, ok := objectStatus.(*status.PersistentVolumeClaimStatus)
	if ok {
		p.Status = *sta
	}
	return ok
}

func (p *PersistentVolumeClaim) GetStatus() ApiObjectStatus {
	return &p.Status
}

func (p *PersistentVolumeClaim) Info() {
	fmt.Printf("%-10s\t%-40s\t%10s\t%-10s\t%-10s\n", "NAME", "UID", "VOLUME", "SCLASS", "PHASE")
	fmt.Printf("%-10s\t%-40s\t%10s\t%-10s\t%-10s\n", p.Metadata.Name, p.Metadata.Uid, p.Spec.VolumeName, p.Spec.StorageClassName, p.Status.Phase)
}

type PersistentVolumeClaimList struct {
	ApiVersion string                  `json:"apiVersion,omitempty"`
	Kind       string                  `json:"kind,omitempty"`
	Items      []PersistentVolumeClaim `json:"items"`
}

func (p *PersistentVolumeClaimList) JsonUnmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, p)
}

func (p *PersistentVolumeClaimList) JsonMarshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PersistentVolumeClaimList) AppendItems(objects []string) error {
	for _, object := range objects {
		var pvc PersistentVolumeClaim
		err := pvc.JsonUnmarshal([]byte(object))
		if err != nil {
			return err
		}
		p.Items = append(p.Items, pvc)
	}
	return nil
}

func (p *PersistentVolumeClaimList) GetItems() []ApiObject {
	var items []ApiObject
	items = make([]ApiObject, 0)
	for _, item := range p.Items {
		items = append(items, &item)
	}
	return items
}

func (p *PersistentVolumeClaimList) Info() {
	fmt.Printf("%-10s\t%-40s\t%10s\t%-10s\t%-10s\n", "NAME", "UID", "VOLUME", "SCLASS", "PHASE")
	for _, item := range p.Items {
		fmt.Printf("%-10s\t%-40s\t%10s\t%-10s\t%-10s\n", item.Metadata.Name, item.Metadata.Uid, item.Spec.VolumeName, item.Spec.StorageClassName, item.Status.Phase)
	}
}
