package heartbeat

import (
	"MiniK8S/pkg/api/config"
	apitypes "MiniK8S/pkg/api/types"
	"MiniK8S/pkg/api/watch"
	"MiniK8S/pkg/apiClient/listWatcher"
	"context"
	"github.com/google/uuid"
	"time"
)

type HeartbeatReceiver struct {
	times         map[uuid.UUID]time.Time
	HbListWatcher listWatcher.ListerWatcher
}

func NewHeartbeatReceiver() *HeartbeatReceiver {
	return &HeartbeatReceiver{
		times: make(map[uuid.UUID]time.Time),
	}
}

func (hbr *HeartbeatReceiver) Run(ctx context.Context, cancel context.CancelFunc) {

	//run list watch
	go hbr.WatchList(ctx, hbr.HbListWatcher)
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
			time.Sleep(config.HeartbeatCheckInterval)
		}
	}
}

func (hbr *HeartbeatReceiver) WatchList(ctx context.Context, listWatcher listWatcher.ListerWatcher) {
	podList, err := hbr.HbListWatcher.List(config.ListOptions{
		Kind:            string(apitypes.HeartbeatObjectType),
		APIVersion:      "",
		LabelSelector:   "",
		FieldSelector:   "",
		Watch:           false,
		ResourceVersion: "",
		TimeoutSeconds:  nil,
	})
	if err != nil {
		panic(err)
	}
	list := podList.GetIApiObjectArr()
	for _, object := range list {
		hb := object.(*config.Heartbeat)
		hbr.times[hb.Metadata.Uid], _ = time.Parse(time.DateTime, hb.Metadata.CreationTimestamp)
	}
	w, err := listWatcher.Watch(config.ListOptions{
		Kind:  string(apitypes.HeartbeatObjectType),
		Watch: true,
	})
	if err != nil {
		panic(err)
	}
	err = hbr.handleWatch(w, ctx)
	if err != nil {
		panic(err)
		return
	}

}

func (hbr *HeartbeatReceiver) handleWatch(w watch.Interface, ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			panic("heartbeat receiver stopped")
			return nil
		case event := <-w.ResultChan():
			hb := event.Object.(*config.Heartbeat)
			ti, _ := time.Parse(time.RFC3339Nano, hb.Metadata.CreationTimestamp)
			switch event.Type {
			case watch.Added:
				hbr.times[hb.Metadata.Uid] = ti
			case watch.Modified:
				hbr.times[hb.Metadata.Uid] = ti
			case watch.Deleted:
				delete(hbr.times, hb.Metadata.Uid)
			case watch.Error:
				panic("heartbeat receiver watch error")
			case watch.Bookmark:
				panic("heartbeat receiver watch bookmark")
			default:
				panic("should not happen")
			}
		}
	}

	return nil
}
