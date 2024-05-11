package pod

import (
	"MiniK8S/pkg/api/config/podConfig"
	"github.com/google/uuid"
)

type containerNameToContainerID map[string]string

type PodManager struct {
	id2Pod     map[uuid.UUID]*podConfig.PodConfig
	name2Pod   map[string]*podConfig.PodConfig
	containers map[string]containerNameToContainerID
	pods       []*podConfig.PodConfig
}

func NewPodManager() *PodManager {
	return &PodManager{
		id2Pod:     map[uuid.UUID]*podConfig.PodConfig{},
		name2Pod:   map[string]*podConfig.PodConfig{},
		containers: map[string]containerNameToContainerID{},
		pods:       []*podConfig.PodConfig{},
	}
}

func (p *PodManager) GetPodById(id uuid.UUID) *podConfig.PodConfig {
	return p.id2Pod[id]
}

func (p *PodManager) GetPodByName(name string) *podConfig.PodConfig {
	return p.name2Pod[name]
}

func (p *PodManager) AddPod(id uuid.UUID, name string, config *podConfig.PodConfig) {
	p.id2Pod[id] = config
	p.name2Pod[name] = config
}

func (p *PodManager) DeletePodById(id uuid.UUID) {
	pod := p.GetPodById(id)
	name := pod.Meta.Name
	delete(p.name2Pod, name)
	delete(p.id2Pod, id)
}

func (p *PodManager) GetPods() []*podConfig.PodConfig {
	return p.pods
}

func (p *PodManager) AddContainer(name string, id string) {
	//p.containers[name] = append(p.containers[name], id)
}
func (p *PodManager) DeletePodByName(name string) {
	pod := p.GetPodByName(name)
	id := pod.Meta.Uid
	delete(p.name2Pod, name)
	delete(p.id2Pod, id)
}

func (p *PodManager) MakePodName(name string, namespace string) string {
	return namespace + "_" + name
}
