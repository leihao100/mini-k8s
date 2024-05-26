package heartbeat

import (
	"MiniK8S/pkg/api/config"
	"context"
	"github.com/google/uuid"
	"time"
)

type HeartbeatReceiver struct {
	times map[uuid.UUID]time.Time
}

func NewHeartbeatReceiver() *HeartbeatReceiver {
	return &HeartbeatReceiver{
		times: make(map[uuid.UUID]time.Time),
	}
}

func (hbr *HeartbeatReceiver) Run(ctx context.Context, cancel context.CancelFunc) {

	//run check
	go func() {
		defer cancel()
		hbr.Check(ctx)
	}()
}

func (hbr *HeartbeatReceiver) Check(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, t := range hbr.times {
				if time.Since(t) > config.HeartbeatTimeoutInterval {
					//handle timeout

				}
			}
			time.Sleep(time.Duration(config.HeartbeatCheckInterval))
		}

	}
}
