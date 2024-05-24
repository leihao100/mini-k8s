package main

import (
	apiserver "MiniK8S/pkg/apiServer"
	"MiniK8S/pkg/kubelet"
	"context"
	"time"
)

func main() {
	go func() {
		api := apiserver.NewApiServer()
		context, cancel := context.WithCancel(context.Background())
		api.Run(cancel)
		<-context.Done()
	}()
	time.Sleep(3 * time.Second)
	k := kubelet.Kubelet{}
	k.Run()
	k.SendMessage()
}
