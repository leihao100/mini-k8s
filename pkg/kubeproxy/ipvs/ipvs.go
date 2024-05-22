package ipvsManager

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/kubeproxy/ipInterface"
	"github.com/cloudflare/ipvs"
	"github.com/cloudflare/ipvs/netmask"
	"net/netip"
)

// IPSet :gets from k8s' source code,may be in use
type IPSet struct {
	Name string
	// SetType specifies the ipset type.
	SetType string
	// HashFamily specifies the protocol family of the IP addresses to be stored in the set.
	// The default is inet, i.e IPv4.  If users want to use IPv6, they should specify inet6.
	HashFamily string
	// HashSize specifies the hash table size of ipset.
	HashSize int
	// MaxElem specifies the max element number of ipset.
	MaxElem int
	// PortRange specifies the port range of bitmap:port type ipset.
	PortRange string
	// comment message for ipset
	Comment string
	//Map maps the service and its ips
	Map map[string]string
}

type IPVSManager struct {
	Client *ipvs.Client
}

func GetIPVS() ipInterface.IP {
	ip, _ := ipvs.New()
	return &IPVSManager{
		Client: &ip,
	}
}

func (im *IPVSManager) New() {
	client, err := ipvs.New()
	if err != nil {
		return
	}
	im.Client = &client
}

func (im *IPVSManager) AddService(service *config.Service) {
	client := *im.Client
	addr, err := netip.ParseAddr(service.Spec.ClusterIP)
	if err != nil {
		return
	}
	for _, port := range service.Spec.Ports {
		client.CreateService(ipvs.Service{
			Address:   addr,
			Port:      uint16(port.Port),
			Netmask:   netmask.Mask{},
			Scheduler: defaultScheduler,
			Family:    defaultAddressFamily,
			Protocol:  defaultProtocol,
		})
	}

}

func (im *IPVSManager) RemoveService(service *config.Service) {
	client := *im.Client
	addr, err := netip.ParseAddr(service.Spec.ClusterIP)
	if err != nil {
		return
	}
	services, err := client.Services()
	if err != nil {
		return
	}
	ports := make([]uint16, 0)
	for _, port := range service.Spec.Ports {
		ports = append(ports, uint16(port.Port))
	}
	for _, port := range ports {
		for _, service := range services {
			if service.Port == port && service.Address == addr {
				client.RemoveService(service.Service)
			}
		}
	}
}

func (im *IPVSManager) AddPodToService(serviceArg *config.Service, pod *config.Pod) {
	client := *im.Client
	addr, err := netip.ParseAddr(serviceArg.Spec.ClusterIP)
	if err != nil {
		return
	}
	services, err := client.Services()
	if err != nil {
		return
	}
	ports := make([]config.ServicePort, 0)
	for _, port := range serviceArg.Spec.Ports {
		ports = append(ports, port)
	}
	for _, port := range ports {
		for _, service := range services {
			//todo 这里用等号可能有问题
			if service.Port == uint16(port.Port) && service.Address == addr {
				//now we have found this
				podIP, _ := netip.ParseAddr(pod.Status.PodIP)
				client.CreateDestination(service.Service, ipvs.Destination{
					Address: podIP,
					//tobe finish
					FwdMethod:      0,
					Weight:         1,
					UpperThreshold: 0,
					LowerThreshold: 0,
					Port:           uint16(port.TargetPort),
					Family:         defaultAddressFamily,
					TunnelType:     0,
					TunnelPort:     0,
					TunnelFlags:    0,
				})
			}
		}
	}
}

func (im *IPVSManager) RemovePodFromService(serviceArg *config.Service, pod *config.Pod) {
	client := *im.Client
	addr, err := netip.ParseAddr(serviceArg.Spec.ClusterIP)
	if err != nil {
		return
	}
	podAddr, err := netip.ParseAddr(pod.Status.PodIP)
	services, err := client.Services()
	if err != nil {
		return
	}
	ports := make([]config.ServicePort, 0)
	for _, port := range serviceArg.Spec.Ports {
		ports = append(ports, port)
	}
	for _, port := range ports {
		for _, service := range services {
			if service.Port == uint16(port.Port) && service.Address == addr {
				//now we have found this
				//then go around its destinations and delete the destination whose ip equals to podIP a
				destinations, _ := client.Destinations(service.Service)
				for _, destination := range destinations {
					if destination.Address == podAddr {
						client.RemoveDestination(service.Service, destination.Destination)
						break
					}
				}
				break
			}
		}
	}
}
