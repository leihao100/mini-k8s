package heartbeat

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"context"
	"github.com/google/uuid"
	"log"
	"time"
)

type HeartBeatSender struct {
	Client *apiClient.Client
	NodeID uuid.UUID
	Hb     *config.Heartbeat
}

func NewHbSender(nodeID uuid.UUID) *HeartBeatSender {
	cli := apiClient.NewRESTClient(types.HeartbeatObjectType)
	hb := &config.Heartbeat{
		Kind: "heartbeat",
		Metadata: meta.ObjectMeta{
			CreationTimestamp: time.Now().String(),
		},
	}
	return &HeartBeatSender{
		NodeID: nodeID,
		Client: cli,
		Hb:     hb,
	}
}

func (hbs *HeartBeatSender) Run(ctx context.Context, cancel context.CancelFunc) {

	go func() {
		defer cancel()
		defer log.Printf("[HeartbeatSender] finished\n")
		hbs.periodicallySendHeartbeat(ctx)
	}()
}

func (hbs *HeartBeatSender) SendHeartbeat() error {

	//hbItem, err := s.heartbeatClient.Get(s.hb.UID)
	//if err != nil {
	//	log.Printf("[updateAndSendHeartbeat] node %v get heartbeat info failed\n", s.nodeUID)
	//	return err
	//}
	//s.hb = hbItem.(*core.Heartbeat)
	hbs.Hb.Metadata.Uid = uuid.New()
	hbs.Hb.Metadata.CreationTimestamp = time.Now().String()
	url := hbs.Client.BuildURL(apiClient.Create)
	buf, _ := hbs.Hb.JsonMarshal()
	res := hbs.Client.Put(url, buf)
	defer func() {
		err := res.Close()
		if err != nil {
			panic("send heartbeat func close http fail")
		}
	}()

	//todo error handle and version update
	return nil

}

func (hbs *HeartBeatSender) periodicallySendHeartbeat(ctx context.Context) {
	// go wait.UntilWithContext(ctx, rsc.worker, time.Second)
	for {
		select {
		case <-ctx.Done():
			log.Printf("[periodicallySendHeartbeat] ctx.Done() received, heartbeat sender exit\n")
			return
		default:
			//send heartbeat by sending heartbeat object to ApiServer
			err := hbs.SendHeartbeat()
			if err != nil {
				panic(err)
			}
			time.Sleep(config.HeartbeatSendInterval)
		}
	}
}
