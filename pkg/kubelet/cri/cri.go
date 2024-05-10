package cri

import (
	"MiniK8S/pkg/util/config/containerConfig"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type Client interface {
	CreateContainer(Config containerConfig.ContainerConfig, name string) (container.CreateResponse, error)
	StartContainer(id string) (bool, error)
	StopContainer(id string) (bool, error)
	ContainerStatus(id string) (bool, int, error)
	RemoveContainer(id string) error
	Close() error
	ListContainers() []types.Container
}
