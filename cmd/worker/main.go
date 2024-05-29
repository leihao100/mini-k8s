package main

import (
	"MiniK8S/pkg/kubelet"
	"MiniK8S/pkg/node"
	"MiniK8S/pkg/node/heartbeat"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	workNode := node.CreateWorkerNode()
	hbSender := heartbeat.NewHbSender(workNode.GetNode().GetUID())
	hbSender.Run(ctx, cancel)
	kubelet.NewKubelet(*workNode.GetNode())

}
