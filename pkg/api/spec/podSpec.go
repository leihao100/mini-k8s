package spec

import (
	"MiniK8S/pkg/api/config/containerConfig"
)

type PodSpec struct {
	Containers     []containerConfig.ContainerConfig
	InitContainers []containerConfig.ContainerConfig
	NodeName       string
	ExposedPorts   []string
	Volumes        []string
	BindPorts      map[string]string
}
