package spec

import (
	"MiniK8S/pkg/util/config/containerConfig"
)

type PodSpec struct {
	Containers     []containerConfig.ContainerConfig
	InitContainers []containerConfig.ContainerConfig
	NodeName       string
	ExposedPorts   []string
	Volumes        []string
	BindPorts      map[string]string
}
