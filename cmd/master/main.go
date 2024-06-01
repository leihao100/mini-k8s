package main

import (
	apiserver "MiniK8S/pkg/apiServer"
	"MiniK8S/pkg/node"
	"MiniK8S/pkg/node/heartbeat"
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

	server := apiserver.NewApiServer()
	server.Run(cancel)
	fmt.Println("server成功运行")
	time.Sleep(5 * time.Second)
	node.CreateMasterNode()
	heartbeatRecevier := heartbeat.NewHeartbeatReceiver()
	heartbeatRecevier.Run(ctx, cancel)
	<-ctx.Done()
}
