package main

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"fmt"
)

func mmain() {
	//k := kubelet.Kubelet{}
	//k.Run()
	pod := config.Pod{
		ApiVersion: "",
		Kind:       "pod",
		Metadata: meta.ObjectMeta{
			Name:      "try",
			Namespace: "try",
		},
		Spec:   config.PodSpec{},
		Status: status.PodStatus{},
	}
	fmt.Println(pod.Metadata.Uid)
}
