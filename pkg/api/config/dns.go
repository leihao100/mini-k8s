package config

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

type DNS struct {
	ApiVersion string           `json:"apiversion,omitempty"`
	Kind       string           `json:"kind,omitempty"`
	Metadata   meta.ObjectMeta  `json:"metadata,omitempty"`
	Spec       DNSSpec          `json:"spec,omitempty"`
	Status     status.DNSStatus `json:"status,omitempty"`
}

type DNSSpec struct {
	HostName string    `json:"hostname,omitempty"`
	HostPort string    `json:"hostport,omitempty"`
	Path     []DNSPath `json:"path,omitempty"`
}
type DNSPath struct {
	ServiceName string `json:"serviceName,omitempty"`
	ClusterIP   string `json:"clusterIP,omitempty"`
	ClusterPort string `json:"clusterPort,omitempty"`
	ClusterPath string `json:"clusterPath,omitempty"`
}

func (d *DNS) JsonMarshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *DNS) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &d)
}

func (d *DNS) SetUID(uid uuid.UUID) {
	d.Metadata.Uid = uid
}

func (d *DNS) GetUID() uuid.UUID {
	return d.Metadata.Uid
}

func (d *DNS) GetName() string {
	return d.Metadata.Name
}

func (d *DNS) SetResourceVersion(version int64) {
	d.Metadata.ResourceVersion = strconv.FormatInt(version, 10)
}
func (d *DNS) GetResourceVersion() int64 {
	res, err := strconv.ParseInt(d.Metadata.ResourceVersion, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	return res
}
func (d *DNS) JsonUnmarshalStatus(data []byte) error {
	return json.Unmarshal(data, &(d.Status))
}

func (d *DNS) JsonMarshalStatus() ([]byte, error) {
	return json.Marshal(d.Status)
}
func (d *DNS) SetStatus(s ApiObjectStatus) bool {
	status, ok := s.(*status.DNSStatus)
	if ok {
		d.Status = *status
	}
	return true
}
func (d *DNS) GetStatus() ApiObjectStatus {
	return &d.Status
}
func (d *DNS) Info() {
	fmt.Printf("%-10s\t%-40s\t%-20s\n", "NAME", "UID", "HOSTNAME")
	fmt.Printf("%-10s\t%-40s\t%-20s\n", d.Metadata.Name, d.Metadata.Uid, d.Spec.HostName)
}

type DNSList struct {
	ApiVersion      string `json:"apiVersion,omitempty"`
	Kind            string `json:"kind,omitempty"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
	Continue        string `json:"continue,omitempty"`
	Items           []DNS  `json:"items"`
}

func (d *DNSList) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &d)
}

func (d *DNSList) JsonMarshal() ([]byte, error) {
	return json.Marshal(d)
}
func (d *DNSList) AppendItems(objects []string) error {
	for _, object := range objects {
		ApiObject := &DNS{}
		err := ApiObject.JsonUnmarshal([]byte(object))
		if err != nil {
			return err
		}
		d.Items = append(d.Items, *ApiObject)
	}
	return nil
}

func (d *DNSList) GetItems() []ApiObject {
	var items []ApiObject
	items = make([]ApiObject, 0)
	for _, item := range d.Items {
		items = append(items, &item)
	}
	return items
}
func (d *DNSList) Info() {
	fmt.Printf("%-10s\t%-40s\t%-20s\n", "NAME", "UID", "HOSTNAME")
	for _, item := range d.Items {
		fmt.Printf("%-10s\t%-40s\t%-20s\n", item.Metadata.Name, item.Metadata.Uid, item.Spec.HostName)
	}
}
