package main

import (
	"MiniK8S/pkg/controller/cache"
	"MiniK8S/pkg/kubelet/cadvisor"
	"fmt"
	v1 "github.com/google/cadvisor/info/v1"
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
	cli := cadvisor.NewCAdvisor("http://localhost:9000")
	query := v1.ContainerInfoRequest{
		NumStats: 12,
	}
	res, _ := cli.Inspect(&query)
	for _, stat := range res {
		fmt.Println(stat)
	}
	//dir, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(dir)
	//dns := config.DNS{
	//	ApiVersion: "",
	//	Kind:       "",
	//	Metadata: meta.ObjectMeta{
	//		Uid: uuid.New(),
	//	},
	//	Spec: config.DNSSpec{
	//		HostName: "leihao.com",
	//		HostPort: "80",
	//		Path: []config.DNSPath{
	//			config.DNSPath{
	//				ClusterIP:   "127.0.0.1",
	//				ClusterPath: "/hello",
	//			},
	//		},
	//	},
	//	Status: status.DNSStatus{},
	//}
	//dns1 := config.DNS{
	//	ApiVersion: "",
	//	Kind:       "",
	//	Metadata: meta.ObjectMeta{
	//		Uid: uuid.New(),
	//	},
	//	Spec: config.DNSSpec{
	//		HostName: "lei.com",
	//		HostPort: "80",
	//		Path: []config.DNSPath{
	//			config.DNSPath{
	//				ClusterIP:   "127.0.0.5",
	//				ClusterPath: "/hello",
	//			},
	//		},
	//	},
	//	Status: status.DNSStatus{},
	//}
	//net.GenerateNginxConfig(dns)
	//net.GenerateNginxConfig(dns1)
	//err := net.RemoveNginxConfig(dns)
	//if err != nil {
	//	fmt.Println(err)
	//}
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

	//client, err := ipvs.
	//if err != nil {
	//	return
	//}
	//svc := ipvs.Service{
	//	Address: net.IP("10.0.6.1"),
	//	Port:    uint16(80),
	//	//AddressFamily: 2,
	//	//Netmask:  netmask.Mask{},
	//	Protocol: 6,
	//}
	//if &svc == nil {
	//	fmt.Println("service is nil")
	//}
	//ipvs.AddService(svc)
	//ipvs.AddDestination(svc, ipvs.Destination{
	//	Address: net.IP("10.0.2.3"),
	//	//tobe finish
	//	//FwdMethod:      0,
	//	Weight:         1,
	//	UpperThreshold: 0,
	//	LowerThreshold: 0,
	//	Port:           uint16(80),
	//	//Family:         defaultAddressFamily,
	//	//TunnelType:     0,
	//	//TunnelPort:     0,
	//	//TunnelFlags:    0,
	//})
	//ipvs.Flush()
	//store := cache.NewSimpleStore()
	//li := store.List()
	//for _, i := range li {
	//	fmt.Println(i)
	//}
	//aaa(store)
	//li = store.List()
	//for _, i := range li {
	//	fmt.Println(i)
	//}
	//queue := cache.NewWorkQueue()
	//queue.Add("132")
	//res, ok := queue.Get()
	//if !ok {
	//	fmt.Println(res)
	//} else {
	//	fmt.Println("queue is empty")
	//}
	////go bbb(queue)
	//res, ok = queue.Get()
	//if !ok {
	//	fmt.Println(res)
	//} else {
	//	fmt.Println("queue is empty")
	//}

}

func aaa(st cache.Store) {
	st.Add("123", "132")
}

func bbb(st *cache.WorkQueue) {
	st.Add("123")
}
