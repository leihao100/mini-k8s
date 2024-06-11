package main

import (
	apiserver "MiniK8S/pkg/apiServer"
	"MiniK8S/pkg/controller"
	gpu "MiniK8S/pkg/gpu/server"
	"MiniK8S/pkg/node"
	"MiniK8S/pkg/node/heartbeat"
	"MiniK8S/pkg/scheduler"
	"context"
	"fmt"
	"time"
)

/*
docker run --rm --net=host quay.io/coreos/etcd:v3.5.13 etcdctl --endpoints=http://localhost:2379 member list
*/
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//config.SetEtcdHost()
	server := apiserver.NewApiServer()
	server.Run(cancel)
	fmt.Println("server成功运行")
	time.Sleep(5 * time.Second)
	node.CreateMasterNode("master")
	heartbeatRecevier := heartbeat.NewHeartbeatReceiver()
	heartbeatRecevier.Run(ctx, cancel)
	newScheduler := scheduler.NewScheduler()
	newScheduler.Run(ctx, cancel)
	controllerManager := controller.NewControllerManager()
	controllerManager.Run(ctx, cancel)
	gpuServer := gpu.NewServer()
	gpuServer.Run(ctx, cancel)
	<-ctx.Done()
}
