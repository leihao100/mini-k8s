package main

import (
	"MiniK8S/config"
	"MiniK8S/pkg/kubelet"
	"MiniK8S/pkg/node"
	"MiniK8S/pkg/node/heartbeat"
	"context"
	"flag"
	"fmt"
)

func main() {
	name := flag.String("n", "", "The name to process")
	host := flag.String("h", "", "The hostname to process")
	port := flag.String("p", "", "The port to process")
	flag.Parse()
	if *name == "" {
		fmt.Println("You must specify a node name")
		return
	}
	if *host == "" {
		fmt.Println("You must specify a hostname")
		return
	}
	if *port == "" {
		*port = ":8080"
	} else {
		*port = ":" + *port
	}

	config.SetApiServerHost("http://" + *host)
	config.SetApiServerPort(*port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	workNode := node.CreateWorkerNode(*name)
	hbSender := heartbeat.NewHbSender(workNode.GetNode().GetUID())
	hbSender.Run(ctx, cancel)
	kubelet := kubelet.NewKubelet(*workNode.GetNode())
	kubelet.Run(ctx, cancel)
	<-ctx.Done()
}
