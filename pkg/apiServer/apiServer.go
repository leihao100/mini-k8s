package apiserver

import (
	"context"
)

type ApiServer interface {
	Run(cancel context.CancelFunc)
}

type apiServer struct {
	httpServer httpServer
}
