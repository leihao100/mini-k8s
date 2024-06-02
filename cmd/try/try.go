package main

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"MiniK8S/utils/net"
	"fmt"
	"github.com/google/uuid"
)

func main() {
	//fmt.Println(testutil.GenerateRandomAlphaOnlyString(9))
	//t := time.Now().Format(time.RFC3339Nano)
	//fmt.Println(t)
	////layout := "2024-05-27 11:22:26.2743388"
	//nt, err := time.Parse(time.RFC3339Nano, t)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//time.Sleep(1 * time.Second)
	//fmt.Println(nt)
	//fmt.Println(time.Since(nt))
	//cli := cadvisor.NewCAdvisor("http://localhost:8080")
	//query := v1.ContainerInfoRequest{
	//	NumStats: 12,
	//}
	//res, _ := cli.Inspect(&query)
	//for _, stat := range res {
	//	fmt.Println(stat)
	//}
	//dir, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(dir)
	dns := config.DNS{
		ApiVersion: "",
		Kind:       "",
		Metadata: meta.ObjectMeta{
			Uid: uuid.New(),
		},
		Spec: config.DNSSpec{
			HostName: "leihao.com",
			HostPort: "80",
			Path: []config.DNSPath{
				config.DNSPath{
					ClusterIP:   "127.0.0.1",
					ClusterPath: "/hello",
				},
			},
		},
		Status: status.DNSStatus{},
	}
	dns1 := config.DNS{
		ApiVersion: "",
		Kind:       "",
		Metadata: meta.ObjectMeta{
			Uid: uuid.New(),
		},
		Spec: config.DNSSpec{
			HostName: "lei.com",
			HostPort: "80",
			Path: []config.DNSPath{
				config.DNSPath{
					ClusterIP:   "127.0.0.5",
					ClusterPath: "/hello",
				},
			},
		},
		Status: status.DNSStatus{},
	}
	net.GenerateNginxConfig(dns)
	net.GenerateNginxConfig(dns1)
	err := net.RemoveNginxConfig(dns)
	if err != nil {
		fmt.Println(err)
	}
	//co := config.Container{
	//	Name:         "helloworld",
	//	Cmd:          nil,
	//	Entrypoint:   nil,
	//	Env:          nil,
	//	Image:        "nginx:latest",
	//	Volumes:      nil,
	//	Labels:       nil,
	//	PortBindings: nil,
	//	VolumesFrom:  nil,
	//	Binds:        nil,
	//	NetworkMode:  "",
	//	CPULimit:     0,
	//	MemLimit:     0,
	//}
	//cli, _ := cri.GetClient()
	//cli.CreateContainer(co, "111")

}
