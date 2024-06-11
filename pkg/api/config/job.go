package config

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

type Job struct {
	ApiVersion string           `json:"apiVersion,omitempty"`
	Kind       string           `json:"kind,omitempty"`
	Metadata   meta.ObjectMeta  `json:"metadata,omitempty"`
	Spec       JobSpec          `json:"spec,omitempty"`
	Status     status.JobStatus `json:"status,omitempty"`
}

type JobSpec struct {
	CudaFilePath   string  `json:"cudaFilePath,omitempty"`
	ResultFileName string  `json:"resultFileName,omitempty"`
	ResultFilePath string  `json:"resultFilePath,omitempty"`
	Args           JobArgs `json:"args,omitempty"`
}

type JobArgs struct {
	Mail            *MailRemind `json:"mail,omitempty"`
	NumTasksperNode int         `json:"numTasksperNode,omitempty"`
	CpusPerTask     int         `json:"cpusPerTask,omitempty"`
	GpuResources    int         `json:"gpuResources,omitempty"`
}

type MailRemind struct {
	Type     MailRemindType `json:"type,omitempty"`
	UserName string         `json:"userName,omitempty"`
}
type MailRemindType string

const (
	MailRemindAll   MailRemindType = "all"
	MailRemindBegin MailRemindType = "begin"
	MailRemindEnd   MailRemindType = "end"
	MailRemindFail  MailRemindType = "fail"
)

type JobList struct {
	ApiVersion      string `json:"apiVersion,omitempty"`
	Kind            string `json:"kind,omitempty"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
	Continue        string `json:"continue,omitempty"`
	Items           []Job  `json:"items"`
}

func (j *Job) Info() {
	state := j.Status.State
	if state != "COMPLETED" {
		state = "RUNNING"
	}
	fmt.Printf("%-10s\t%-40s\t%-20s\n", "NAME", "UID", "STATE")
	fmt.Printf("%-10s\t%-40s\t%-20s\n", j.Metadata.Name, j.Metadata.Uid, state)
}

func (j *Job) SetUID(uid uuid.UUID) {
	j.Metadata.Uid = uid
}

func (j *Job) GetUID() uuid.UUID {
	return j.Metadata.Uid
}

func (j *Job) GetName() string {
	return j.Metadata.Name
}

func (j *Job) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &j)
}

func (j *Job) JsonMarshal() ([]byte, error) {
	return json.Marshal(j)
}

func (j *Job) JsonUnmarshalStatus(data []byte) error {
	return json.Unmarshal(data, &(j.Status))
}

func (j *Job) JsonMarshalStatus() ([]byte, error) {
	return json.Marshal(j.Status)
}

func (j *Job) SetStatus(s ApiObjectStatus) bool {
	status, ok := s.(*status.JobStatus)
	if ok {
		j.Status = *status
	}
	return ok
}

func (j *Job) GetStatus() ApiObjectStatus {
	return &j.Status
}

func (j *Job) GetResourceVersion() int64 {
	res, err := strconv.ParseInt(j.Metadata.ResourceVersion, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	return res
}

func (j *Job) SetResourceVersion(version int64) {
	j.Metadata.ResourceVersion = strconv.FormatInt(version, 10)
}

func (j *JobList) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &j)
}

func (j *JobList) JsonMarshal() ([]byte, error) {
	return json.Marshal(j)
}
func (j *JobList) AppendItems(objects []string) error {
	for _, object := range objects {
		ApiObject := &Job{}
		err := ApiObject.JsonUnmarshal([]byte(object))
		if err != nil {
			return err
		}
		j.Items = append(j.Items, *ApiObject)
	}
	return nil
}
func (j *JobList) GetItems() []ApiObject {
	var items []ApiObject
	items = make([]ApiObject, 0)
	for _, item := range j.Items {
		items = append(items, &item)
	}
	return items
}
func (j *JobList) Info() {
	fmt.Printf("%-10s\t%-40s\t%-20s\n", "NAME", "UID", "STATE")
	for _, item := range j.Items {
		state := item.Status.State
		if state != "COMPLETED" {
			state = "RUNNING"
		}
		fmt.Printf("%-10s\t%-40s\t%-20s\n", item.Metadata.Name, item.Metadata.Uid, state)
	}
}
