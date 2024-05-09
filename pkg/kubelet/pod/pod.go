package pod

import "MiniK8S/pkg/util/config"

type podManager struct {
	id2Pod   map[int64]*config.PodConfig
	name2Pod map[string]*config.PodConfig
}

func NewPodManager() *podManager {
	return &podManager{
		id2Pod:   map[int64]*config.PodConfig{},
		name2Pod: map[string]*config.PodConfig{},
	}
}

func (p *podManager) GetPodById(id int64) config.PodConfig {
	return *p.id2Pod[id]
}

func (p *podManager) AddPod(id int64, name string, config *config.PodConfig) {
	p.id2Pod[id] = config
	p.name2Pod[name] = config
}

func (p *podManager) DeletePodById(id int64) {
	pod := p.GetPodById(id)
	name := pod.PodName

	delete(p.name2Pod, name)
	delete(p.id2Pod, id)
}
func (p *podManager) DeletePodByName(name string) {
	//pod := p.GetPod(id)
	//name := pod.PodName
	//
	//delete(p.name2Pod, name)
	//delete(p.id2Pod, id)
}
