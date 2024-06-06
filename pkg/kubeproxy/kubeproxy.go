package kubeproxy

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/selector"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/api/watch"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/apiClient/listwatch"
	"MiniK8S/pkg/kubeproxy/ipInterface"
	iptableManager "MiniK8S/pkg/kubeproxy/iptable"
	"MiniK8S/utils/net"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type KubeProxy struct {
	//kl            *kubelet.Kubelet
	services      map[uuid.UUID]*config.Service
	ipManager     ipInterface.IP
	serviceToPods map[uuid.UUID][]*config.Pod

	serviceClient      *apiClient.Client
	podClient          *apiClient.Client
	dnsClient          *apiClient.Client
	serviceListWatcher listwatch.ListerWatcher
	podListWatcher     listwatch.ListerWatcher
	dnsListWatcher     listwatch.ListerWatcher
}

func NewKubeProxy() *KubeProxy {
	return &KubeProxy{
		//kl:                 kl,
		services:           make(map[uuid.UUID]*config.Service),
		ipManager:          iptableManager.New(),
		serviceToPods:      make(map[uuid.UUID][]*config.Pod),
		serviceClient:      apiClient.NewRESTClient(types.ServiceObjectType),
		podClient:          apiClient.NewRESTClient(types.PodObjectType),
		dnsClient:          apiClient.NewRESTClient(types.DnsObjectType),
		serviceListWatcher: nil,
		podListWatcher:     nil,
	}
}

func (kp *KubeProxy) Run(ctx context.Context) {
	fmt.Println("[kube-proxy] Starting KubeProxy")
	ctx, cancel := context.WithCancel(ctx)
	kp.podListWatcher = listwatch.NewListWatchFromClient(kp.podClient)
	kp.serviceListWatcher = listwatch.NewListWatchFromClient(kp.serviceClient)
	kp.dnsListWatcher = listwatch.NewListWatchFromClient(kp.dnsClient)
	go kp.PodListWatch(ctx, cancel)
	go kp.ServiceListWatch(ctx, cancel)
	go kp.DnsListWatch(ctx, cancel)
	return
}

func (kp *KubeProxy) PodListWatch(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	_, err := kp.podListWatcher.List(config.ListOptions{
		Kind:            string(types.PodObjectType),
		APIVersion:      "",
		LabelSelector:   "",
		FieldSelector:   "",
		Watch:           false,
		ResourceVersion: "",
		TimeoutSeconds:  nil,
	})
	if err != nil {
		return
	}

	w, err := kp.podListWatcher.Watch(config.ListOptions{
		Kind:            string(types.PodObjectType),
		APIVersion:      "",
		LabelSelector:   "",
		FieldSelector:   "",
		Watch:           true,
		ResourceVersion: "",
		TimeoutSeconds:  nil,
	})
	if err != nil {
		panic("kube-proxy: failed to watch pod list")
	}
	err = kp.HandlePodWatch(w, ctx)
	if err != nil {
		return
	}
	w.Stop()

}

func (kp *KubeProxy) ServiceListWatch(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	_, err := kp.serviceListWatcher.List(config.ListOptions{
		Kind:            string(types.ServiceObjectType),
		APIVersion:      "",
		LabelSelector:   "",
		FieldSelector:   "",
		Watch:           false,
		ResourceVersion: "",
		TimeoutSeconds:  nil,
	})
	if err != nil {
		return
	}

	w, err := kp.serviceListWatcher.Watch(config.ListOptions{
		Kind:            string(types.ServiceObjectType),
		APIVersion:      "",
		LabelSelector:   "",
		FieldSelector:   "",
		Watch:           true,
		ResourceVersion: "",
		TimeoutSeconds:  nil,
	})
	if err != nil {
		panic("kube-proxy: failed to watch pod list")
	}
	err = kp.HandleServiceWatch(w, ctx)
	if err != nil {
		return
	}
	w.Stop()
}

func (kp *KubeProxy) DnsListWatch(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	_, err := kp.dnsListWatcher.List(config.ListOptions{
		Kind:            string(types.DnsObjectType),
		APIVersion:      "",
		LabelSelector:   "",
		FieldSelector:   "",
		Watch:           false,
		ResourceVersion: "",
		TimeoutSeconds:  nil,
	})
	if err != nil {
		return
	}

	w, err := kp.dnsListWatcher.Watch(config.ListOptions{
		Kind:            string(types.DnsObjectType),
		APIVersion:      "",
		LabelSelector:   "",
		FieldSelector:   "",
		Watch:           true,
		ResourceVersion: "",
		TimeoutSeconds:  nil,
	})
	if err != nil {
		panic("kube-proxy: failed to watch pod list")
	}
	err = kp.HandleDnsWatch(w, ctx)
	if err != nil {
		return
	}
	w.Stop()

}

func (kp *KubeProxy) CreateService(service *config.Service) {
	fmt.Println("[kube-proxy] Creating service")
	kp.services[service.Metadata.Uid] = service
	kp.ipManager.AddService(service)
	pods := kp.SelectPod(service)
	kp.serviceToPods[service.Metadata.Uid] = pods
	for _, pod := range pods {
		fmt.Println("[kube-proxy] Adding Pod: ", pod.Metadata.Name, "To Service: ", service.Metadata.Name)
		kp.ipManager.AddPodToService(service, pod)
		service.Spec.Endpoints = append(service.Spec.Endpoints, pod.Status.PodIP)

	}
	url := kp.serviceClient.BuildURL(apiClient.Create)
	buf, _ := service.JsonMarshal()
	kp.serviceClient.Put(url, buf)
}

func (kp *KubeProxy) SelectPod(service *config.Service) []*config.Pod {
	pods, _ := kp.podListWatcher.List(config.ListOptions{
		Kind:       string(types.PodObjectType),
		APIVersion: "",
		Watch:      false,
	})

	var targetPods []*config.Pod
	for _, pod := range pods.GetItems() {
		pd := pod.(*config.Pod)
		if selector.LabelCompare(pd.Metadata.Labels, service.Spec.Selector) {
			targetPods = append(targetPods, pd)
		}
	}
	return targetPods
}

func (kp *KubeProxy) RemoveService(service *config.Service) {
	fmt.Println("[kube-proxy] Removing service")
	delete(kp.services, service.Metadata.Uid)
	//kp.services[service.Metadata.Uid] = nil
	kp.ipManager.RemoveService(service)
	pods := kp.serviceToPods[service.Metadata.Uid]
	for _, pod := range pods {
		kp.ipManager.RemovePodFromService(service, pod)
	}
}

func (kp *KubeProxy) RemovePod(pod *config.Pod) {
	svcs, _ := kp.serviceListWatcher.List(config.ListOptions{})
	//scs := svcs.GetItems()
	for _, service := range svcs.GetItems() {
		svc := service.(*config.Service)
		if selector.LabelCompare(svc.Metadata.Labels, pod.Metadata.Labels) {
			kp.ipManager.RemovePodFromService(svc, pod)

			for i, endpoint := range svc.Spec.Endpoints {
				if endpoint == pod.Status.PodIP {
					svc.Spec.Endpoints = append(svc.Spec.Endpoints[:i], svc.Spec.Endpoints[i+1:]...)
					break
				}
			}
			url := kp.serviceClient.BuildURL(apiClient.Create)
			buf, _ := service.JsonMarshal()
			kp.serviceClient.Put(url, buf)
		}
	}
}

func (kp *KubeProxy) AddPod(pod *config.Pod) {

	for _, service := range kp.services {
		if selector.LabelCompare(service.Metadata.Labels, pod.Metadata.Labels) {
			service.Spec.Endpoints = append(service.Spec.Endpoints, pod.Status.PodIP)
			kp.ipManager.AddPodToService(service, pod)
			url := kp.serviceClient.BuildURL(apiClient.Create)
			buf, _ := service.JsonMarshal()
			kp.serviceClient.Put(url, buf)
		}
	}
}

func (kp *KubeProxy) GetSvc() {

}

func (kp *KubeProxy) CreateDns(dns *config.DNS) {
	fmt.Println("[kube-proxy] Creating DNS")
	svcs, _ := kp.serviceListWatcher.List(config.ListOptions{
		Watch: false,
		Kind:  string(types.ServiceObjectType),
	})
	for _, svc := range svcs.GetItems() {
		s := svc.(*config.Service)
		for _, path := range dns.Spec.Path {
			if strings.EqualFold(s.Metadata.Name, path.ServiceName) {
				path.ClusterIP = s.Spec.ClusterIP
				net.GenerateNginxConfig(*dns)
			}
		}
	}
	net.AddHost(dns)
	//url := kp.dnsClient.BuildURL(apiClient.Create)
	//buf, _ := dns.JsonMarshal()

}

func (kp *KubeProxy) RemoveDns(dns *config.DNS) {
	fmt.Println("[kube-proxy] Removing DNS")
	//net.RemoveNginxConfig(dns)
	net.RemoveNginxConfig(*dns)
	net.RemoveHost(*dns)
}

func (kp *KubeProxy) HandlePodWatch(w watch.Interface, ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-w.ResultChan():
			switch event.Type {
			case watch.Added:
				kp.AddPod(event.Object.(*config.Pod))
			case watch.Modified:
			case watch.Deleted:
				kp.RemovePod(event.Object.(*config.Pod))
			case watch.Error:
				panic("kube-proxy: watch pod error")
			case watch.Bookmark:
			default:
				panic("should never get here")

			}
		}
	}
}

func (kp *KubeProxy) HandleServiceWatch(w watch.Interface, ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-w.ResultChan():
			fmt.Println("[kube-proxy] handleServiceWatch event:")
			switch event.Type {
			case watch.Added:
				kp.CreateService(event.Object.(*config.Service))
			case watch.Modified:
			case watch.Deleted:
				kp.RemoveService(event.Object.(*config.Service))
			case watch.Error:
				panic("kube-proxy: watch svc error")
			case watch.Bookmark:
			default:
				panic("should never get here")

			}
		}
	}
}

func (kp *KubeProxy) HandleDnsWatch(w watch.Interface, ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-w.ResultChan():
			switch event.Type {
			case watch.Added:
				kp.CreateDns(event.Object.(*config.DNS))
			case watch.Modified:
			case watch.Deleted:
				kp.RemoveDns(event.Object.(*config.DNS))
			case watch.Error:
				panic("kube-proxy: watch Dns error")
			case watch.Bookmark:
			default:
				panic("should never get here")

			}
		}
	}
}
