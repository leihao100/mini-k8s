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
)

type Time = time.Time
type Quantity string
