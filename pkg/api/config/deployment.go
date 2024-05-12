package config

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
)

type Deployment struct {
	ApiVersion string                  `json:"apiVersion,omitempty"`
	Kind       string                  `json:"kind,omitempty"`
	Metadata   meta.ObjectMeta         `json:"metadata,omitempty"`
	Spec       DeploymentSpec          `json:"spec,omitempty"`
	Status     status.DeploymentStatus `json:"status,omitempty"`
}

type DeploymentSpec struct {
}
