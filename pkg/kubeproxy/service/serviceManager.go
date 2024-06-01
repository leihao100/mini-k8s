package service

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/kubeproxy/ipInterface"
	ipvsManager "MiniK8S/pkg/kubeproxy/ipvs"
	"github.com/google/uuid"
)

type ServiceManager struct {
	ServiceMap map[uuid.UUID]*config.Service
	IPManager  ipInterface.IP
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap: make(map[uuid.UUID]*config.Service),
		IPManager:  ipvsManager.GetIPVS(),
	}
}

func (sm *ServiceManager) Run() {

}

func (sm *ServiceManager) AddService(service *config.Service) {
	sm.ServiceMap[service.Metadata.Uid] = service
	sm.IPManager.AddService(service)
}
