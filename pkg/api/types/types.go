package types

import "time"

type ApiObjectType string

const (
	ErrorObjectType                   ApiObjectType = "Error"
	PodObjectType                     ApiObjectType = "Pod"
	ServiceObjectType                 ApiObjectType = "Service"
	ReplicasetObjectType              ApiObjectType = "ReplicaSet"
	HorizontalPodAutoscalerObjectType ApiObjectType = "HorizontalPodAutoscaler"
	NodeObjectType                    ApiObjectType = "Node"
	JobObjectType                     ApiObjectType = "Job"
	HeartbeatObjectType               ApiObjectType = "Heartbeat"
	FuncTemplateObjectType            ApiObjectType = "Func"
	DnsObjectType                     ApiObjectType = "DNS"
	DeploymentObjectType              ApiObjectType = "deployment"
	StorageClassObjectType            ApiObjectType = "StorageClass"
	PersistentVolumeObjectType        ApiObjectType = "PersistentVolume"
	PersistentVolumeClaimObjectType   ApiObjectType = "PersistentVolumeClaim"
)

type Time = time.Time
type Quantity string

const (
	PhasePending = "Pending"
	PhaseRunning = "Running"
)

type JobState string

const (
	JobPending   JobState = "PENDING"
	JobRunning   JobState = "RUNNING"
	JobFailed    JobState = "FAILED"
	JobCompleted JobState = "COMPLETED"
	JobMissing   JobState = "MISSING"
)
