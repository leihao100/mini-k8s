package heartbeat

import (
	"MiniK8S/pkg/api/config"
	apitypes "MiniK8S/pkg/api/types"
	"MiniK8S/pkg/api/watch"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/apiClient/listwatch"
	"context"
	"time"
)

type HeartbeatReceiver struct {
	times         map[string]time.Time
	HbListWatcher listwatch.ListerWatcher
	nodeClient    *apiClient.Client
}

func NewHeartbeatReceiver() *HeartbeatReceiver {
	return &HeartbeatReceiver{
		times:         make(map[string]time.Time),
		HbListWatcher: listwatch.NewListWatchFromClient(apiClient.NewRESTClient(apitypes.HeartbeatObjectType)),
		nodeClient:    apiClient.NewRESTClient(apitypes.NodeObjectType),
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
			for name, t := range hbr.times {
				if time.Since(t) > config.HeartbeatTimeoutInterval {
					//handle timeout
					url := hbr.nodeClient.BuildURL(apiClient.Delete) + "/" + name
					hbr.nodeClient.Delete(url, nil)

				}
			}
			time.Sleep(config.HeartbeatCheckInterval)
		}
	}
}

func (hbr *HeartbeatReceiver) WatchList(ctx context.Context, listWatcher listwatch.ListerWatcher) {
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
	list := podList.GetItems()
	for _, object := range list {
		hb := object.(*config.Heartbeat)
		hbr.times[hb.Metadata.Name], _ = time.Parse(time.DateTime, hb.Metadata.CreationTimestamp)
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
				hbr.times[hb.Metadata.Name] = ti
			case watch.Modified:
				hbr.times[hb.Metadata.Name] = ti
			case watch.Deleted:
				delete(hbr.times, hb.Metadata.Name)
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
