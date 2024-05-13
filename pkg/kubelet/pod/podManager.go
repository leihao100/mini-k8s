package pod

import (
	"MiniK8S/pkg/api/config"
	"fmt"

	"github.com/google/uuid"
)

type containerNameToContainerID map[string]string

type PodManager struct {
	id2Pod     map[uuid.UUID]*config.Pod
	name2Pod   map[string]*config.Pod
	containers map[uuid.UUID]containerNameToContainerID
	pods       []*config.Pod
}

func NewPodManager() *PodManager {
	return &PodManager{
		id2Pod:     map[uuid.UUID]*config.Pod{},
		name2Pod:   map[string]*config.Pod{},
		containers: map[uuid.UUID]containerNameToContainerID{},
		pods:       []*config.Pod{},
	}
}

func (p *PodManager) GetPodById(id uuid.UUID) *config.Pod {
	return p.id2Pod[id]
}

func (p *PodManager) GetPodByName(name string) *config.Pod {
	return p.name2Pod[name]
}

func (p *PodManager) AddPod(id uuid.UUID, name string, config *config.Pod) {
	p.id2Pod[id] = config
	p.name2Pod[name] = config
	p.containers[id] = containerNameToContainerID{}
	p.pods = append(p.pods, config)
}

func (p *PodManager) DeletePodById(id uuid.UUID) {
	pod := p.GetPodById(id)
	name := pod.Metadata.Name
	delete(p.name2Pod, name)
	delete(p.id2Pod, id)
	delete(p.containers, id)
}

func (p *PodManager) GetPods() []*config.Pod {
	return p.pods
}

func (p *PodManager) AddContainer(podID uuid.UUID, containerName string, id string) {
	containerMap, _ := p.containers[podID]
	if containerMap == nil {
		fmt.Println("container map nil")
	}
	containerMap[containerName] = id
	//p.containers[podID][containerName] = id
}
func (p *PodManager) DeletePodByName(name string) {
	pod := p.GetPodByName(name)
	id := pod.Metadata.Uid
	delete(p.name2Pod, name)
	delete(p.id2Pod, id)
	delete(p.containers, uuid.UUID{})
}

func (p *PodManager) MakePodName(pod *config.Pod) string {
	return pod.Metadata.Namespace + "_" + pod.Metadata.Name
}
