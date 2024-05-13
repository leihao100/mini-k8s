package service

import (
	"MiniK8S/pkg/api/config"
	ipvsManager "MiniK8S/pkg/kubeproxy/ipvs"
	"github.com/google/uuid"
)

type ServiceManager struct {
	ServiceMap  map[uuid.UUID]*config.Service
	IPVSManager ipvsManager.IPVSManager
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap:  make(map[uuid.UUID]*config.Service),
		IPVSManager: ipvsManager.IPVSManager{},
	}
}

func (sm *ServiceManager) Run() {
	sm.IPVSManager.New()
}

func (sm *ServiceManager) AddService(service *config.Service) {
	sm.ServiceMap[service.Metadata.Uid] = service
}
