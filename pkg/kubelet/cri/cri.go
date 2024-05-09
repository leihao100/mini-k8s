package cri

import (
	"MiniK8S/pkg/util/config"
	"github.com/docker/docker/api/types"
)

type Client interface {
	CreateContainer(Config config.ContainerConfig, name string) (bool, error)
	StartContainer(id string) (bool, error)
	StopContainer(id string) (bool, error)
	ContainerStatus(id string) (bool, int, error)
	RemoveContainer(id string) error
	Close() error
	ListContainers() []types.Container
}
