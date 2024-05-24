package main

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/meta"

	"MiniK8S/pkg/api/status"
	"MiniK8S/pkg/kubelet"
)

func ain() {
	co := config.Container{
		Name:         "helloworld",
		Cmd:          nil,
		Entrypoint:   nil,
		Env:          nil,
		Image:        "hello-world:latest",
		Volumes:      nil,
		Labels:       nil,
		PortBindings: nil,
		VolumesFrom:  nil,
		Binds:        nil,
		NetworkMode:  "",
		CPULimit:     0,
		MemLimit:     0,
	}
	co1 := config.Container{
		Name:         "nginx",
		Cmd:          nil,
		Entrypoint:   nil,
		Env:          nil,
		Image:        "nginx:latest",
		Volumes:      nil,
		Labels:       nil,
		PortBindings: nil,
		VolumesFrom:  nil,
		Binds:        nil,
		NetworkMode:  "",
		CPULimit:     0,
		MemLimit:     0,
	}
	var containers []config.Container
	containers = append(containers, co)
	containers = append(containers, co1)
	pod := config.Pod{
		ApiVersion: "",
		Kind:       "pod",
		Metadata: meta.ObjectMeta{
			Name:      "try",
			Namespace: "try",
		},
		Spec: config.PodSpec{
			Containers: containers,
		},
		Status: status.PodStatus{},
	}

	k := kubelet.Kubelet{}
	k.Run()
	k.MakePod(&pod)

}
