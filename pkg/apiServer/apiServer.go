package apiserver

import (
	"MiniK8S/config"
	"MiniK8S/pkg/etcd"
	"context"
	"fmt"
)

type ApiServer interface {
	Run(cancel context.CancelFunc)
}

type apiServer struct {
	httpServer HttpServer
}

func NewApiServer() ApiServer {
	return &apiServer{
		httpServer: NewHttpServer(),
	}
}

func (a apiServer) Run(cancel context.CancelFunc) {
	fmt.Printf("[apiServer] apiServer start\n")
	defer fmt.Printf("[apiServer] apiServer start finish")
	etcd.Init()
	a.httpServer.BindHandlers()

	go func() {
		defer cancel()
		defer etcd.Close()
		fmt.Printf("[apiServer] httpServer start\n")
		err := a.httpServer.Run(config.ApiServerPort())
		if err != nil {
			fmt.Printf("[apiServer] httpServer start failed\n")
		}
	}()
}
