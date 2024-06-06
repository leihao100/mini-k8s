package iptableManager

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/utils/net"
	"fmt"
	"log"
	"slices"
	"strconv"

	"github.com/coreos/go-iptables/iptables"
)

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

type IPtableManager struct {
	Client         *iptables.IPTables
	Service2Ports  map[string][]Ports
	Service2Server map[string][]string
}

type Ports struct {
	Port       int32
	TargetPort int32
}

func removeElement(array []string, value string) []string {
	index := slices.Index(array, value)
	if index != -1 {
		newArray := slices.Delete(array, index, index+1)
		return newArray
	}
	return array
}

func New() *IPtableManager {
	client, err := iptables.New()
	if err != nil {
		fmt.Printf("[kubeproxy] Error creating iptables handle: %v\n", err)
		return nil
	}
	return &IPtableManager{
		Client:         client,
		Service2Ports:  make(map[string][]Ports),
		Service2Server: make(map[string][]string),
	}
}

func (im *IPtableManager) AddRules(clusterIP string, serverIP string, clusterPort string, serverPort string, probility string) {
	client := im.Client
	err := client.AppendUnique("nat", "PREROUTING", "-d", clusterIP, "-p", "tcp", "--dport", clusterPort, "-m", "statistic", "--mode", "random", "--probability", probility, "-j", "DNAT", "--to-destination", serverIP+":"+serverPort)
	if err != nil {
		log.Fatalf("Error adding PREROUTING rule: %v", err)
	}
	// Add the second rule
	err = client.AppendUnique("nat", "POSTROUTING", "-d", serverIP, "-p", "tcp", "--dport", serverPort, "-j", "MASQUERADE")
	if err != nil {
		log.Fatalf("Error adding POSTROUTING rule: %v", err)
	}
	// Add the third rule
	err = client.AppendUnique("nat", "OUTPUT", "-d", clusterIP, "-p", "tcp", "--dport", clusterPort, "-m", "statistic", "--mode", "random", "--probability", probility, "-j", "DNAT", "--to-destination", serverIP+":"+serverPort)
	if err != nil {
		log.Fatalf("Error adding OUTPUT rule: %v", err)
	}
}

func (im *IPtableManager) DeleteRules(clusterIP string, serverIP string, clusterPort string, serverPort string, probility string) {
	client := im.Client
	err := client.Delete("nat", "PREROUTING", "-d", clusterIP, "-p", "tcp", "--dport", clusterPort, "-m", "statistic", "--mode", "random", "--probability", probility, "-j", "DNAT", "--to-destination", serverIP+":"+serverPort)
	if err != nil {
		log.Fatalf("Error deleting PREROUTING rule: %v", err)
	}

	// Delete the second rule
	err = client.Delete("nat", "POSTROUTING", "-d", serverIP, "-p", "tcp", "--dport", serverPort, "-j", "MASQUERADE")
	if err != nil {
		log.Fatalf("Error deleting POSTROUTING rule: %v", err)
	}

	// Delete the third rule
	err = client.Delete("nat", "OUTPUT", "-d", clusterIP, "-p", "tcp", "--dport", clusterPort, "-m", "statistic", "--mode", "random", "--probability", probility, "-j", "DNAT", "--to-destination", serverIP+":"+serverPort)
	if err != nil {
		log.Fatalf("Error deleting OUTPUT rule: %v", err)
	}
}

func (im *IPtableManager) AddRulesNodePort(clusterIP string, serverIP string, clusterPort string, serverPort string, probility string) {
	client := im.Client
	err := client.AppendUnique("nat", "PREROUTING", "-p", "tcp", "--dport", clusterPort, "-m", "statistic", "--mode", "random", "--probability", probility, "-j", "DNAT", "--to-destination", serverIP+":"+serverPort)
	if err != nil {
		log.Fatalf("Error adding PREROUTING rule: %v", err)
	}
	// Add the second rule
	err = client.AppendUnique("nat", "POSTROUTING", "-d", serverIP, "-p", "tcp", "--dport", serverPort, "-j", "SNAT", "--to-source", clusterIP)
	if err != nil {
		log.Fatalf("Error adding POSTROUTING rule: %v", err)
	}
}

func (im *IPtableManager) DeleteRulesNodePort(clusterIP string, serverIP string, clusterPort string, serverPort string, probility string) {
	client := im.Client
	err := client.Delete("nat", "PREROUTING", "-p", "tcp", "--dport", clusterPort, "-m", "statistic", "--mode", "random", "--probability", probility, "-j", "DNAT", "--to-destination", serverIP+":"+serverPort)
	if err != nil {
		log.Fatalf("Error deleting PREROUTING rule: %v", err)
	}

	// Delete the second rule
	err = client.Delete("nat", "POSTROUTING", "-d", serverIP, "-p", "tcp", "--dport", serverPort, "-j", "SNAT", "--to-source", clusterIP)
	if err != nil {
		log.Fatalf("Error deleting POSTROUTING rule: %v", err)
	}
}

func (im *IPtableManager) AddService(service *config.Service) {
	if service.Spec.Type == "NodePort" {
		service.Spec.ClusterIP, _ = net.GetLocalIP()
	}
	cli := apiClient.NewRESTClient(types.ServiceObjectType)
	url := cli.BuildURL(apiClient.Create)
	buf, _ := service.JsonMarshal()
	cli.Put(url, buf)
	for _, port := range service.Spec.Ports {

		im.Service2Ports[service.Metadata.Name] = append(im.Service2Ports[service.Metadata.Name], Ports{
			Port:       port.Port,
			TargetPort: port.TargetPort,
		})
	}
}

func (im *IPtableManager) RemoveService(serviceArg *config.Service) {
	servers := im.Service2Server[serviceArg.Metadata.Name]
	num := len(servers)
	if serviceArg.Spec.Type == "NodePort" {
		for _, Ports := range im.Service2Ports[serviceArg.Metadata.Name] {
			clusterPort := strconv.Itoa(int(Ports.Port))
			serverPort := strconv.Itoa(int(Ports.TargetPort))
			for i, serverIP := range servers {
				probility := strconv.FormatFloat(float64(1/(num-i)), 'f', 4, 64)
				fmt.Printf("Probility: %v\n", probility)
				im.DeleteRulesNodePort(serviceArg.Spec.ClusterIP, serverIP, clusterPort, serverPort, probility)
			}
		}
	} else {
		for _, Ports := range im.Service2Ports[serviceArg.Metadata.Name] {
			clusterPort := strconv.Itoa(int(Ports.Port))
			serverPort := strconv.Itoa(int(Ports.TargetPort))
			for i, serverIP := range servers {
				probility := strconv.FormatFloat(float64(1/(num-i)), 'f', 4, 64)
				fmt.Printf("Probility: %v\n", probility)
				im.DeleteRules(serviceArg.Spec.ClusterIP, serverIP, clusterPort, serverPort, probility)
			}
		}
	}

	delete(im.Service2Ports, serviceArg.Metadata.Name)
	delete(im.Service2Server, serviceArg.Metadata.Name)
}

func (im *IPtableManager) AddPodToService(serviceArg *config.Service, pod *config.Pod) {
	oldServers := im.Service2Server[serviceArg.Metadata.Name]
	im.Service2Server[serviceArg.Metadata.Name] = append(im.Service2Server[serviceArg.Metadata.Name], pod.Status.PodIP)
	newServers := im.Service2Server[serviceArg.Metadata.Name]
	fmt.Printf("Add:\n Old: %v\nNew: %v\n", oldServers, newServers)
	oldNum := len(oldServers)
	newNum := len(newServers)
	if serviceArg.Spec.Type == "NodePort" {
		for _, Ports := range im.Service2Ports[serviceArg.Metadata.Name] {
			clusterPort := strconv.Itoa(int(Ports.Port))
			serverPort := strconv.Itoa(int(Ports.TargetPort))
			for i, serverIP := range oldServers {
				probility := strconv.FormatFloat(float64(1/(oldNum-i)), 'f', 4, 64)
				fmt.Printf("Old Probility: %v\n", probility)
				im.DeleteRulesNodePort(serviceArg.Spec.ClusterIP, serverIP, clusterPort, serverPort, probility)
			}
			for i, serverIP := range newServers {
				probility := strconv.FormatFloat(float64(1/(newNum-i)), 'f', 4, 64)
				fmt.Printf("New Probility: %v\n", probility)
				im.AddRulesNodePort(serviceArg.Spec.ClusterIP, serverIP, clusterPort, serverPort, probility)
			}
		}
	} else {
		for _, Ports := range im.Service2Ports[serviceArg.Metadata.Name] {
			clusterPort := strconv.Itoa(int(Ports.Port))
			serverPort := strconv.Itoa(int(Ports.TargetPort))
			for i, serverIP := range oldServers {
				probility := strconv.FormatFloat(float64(1/(oldNum-i)), 'f', 4, 64)
				fmt.Printf("Old Probility: %v\n", probility)
				im.DeleteRules(serviceArg.Spec.ClusterIP, serverIP, clusterPort, serverPort, probility)
			}
			for i, serverIP := range newServers {
				probility := strconv.FormatFloat(float64(1/(newNum-i)), 'f', 4, 64)
				fmt.Printf("New Probility: %v\n", probility)
				im.AddRules(serviceArg.Spec.ClusterIP, serverIP, clusterPort, serverPort, probility)
			}
		}
	}

}

func (im *IPtableManager) RemovePodFromService(serviceArg *config.Service, pod *config.Pod) {
	oldServers := im.Service2Server[serviceArg.Metadata.Name]
	newServers := removeElement(oldServers, pod.Status.PodIP)
	im.Service2Server[serviceArg.Metadata.Name] = newServers
	fmt.Printf("Remove:\n Old: %v\nNew: %v\n", oldServers, newServers)
	oldNum := len(oldServers)
	newNum := len(newServers)
	if serviceArg.Spec.Type == "NodePort" {
		for _, Ports := range im.Service2Ports[serviceArg.Metadata.Name] {
			clusterPort := strconv.Itoa(int(Ports.Port))
			serverPort := strconv.Itoa(int(Ports.TargetPort))
			for i, serverIP := range oldServers {
				probility := strconv.FormatFloat(float64(1/(oldNum-i)), 'f', 4, 64)
				fmt.Printf("Old Probility: %v\n", probility)
				im.DeleteRulesNodePort(serviceArg.Spec.ClusterIP, serverIP, clusterPort, serverPort, probility)
			}
			for i, serverIP := range newServers {
				probility := strconv.FormatFloat(float64(1/(newNum-i)), 'f', 4, 64)
				fmt.Printf("New Probility: %v\n", probility)
				im.AddRulesNodePort(serviceArg.Spec.ClusterIP, serverIP, clusterPort, serverPort, probility)
			}
		}
	} else {
		for _, Ports := range im.Service2Ports[serviceArg.Metadata.Name] {
			clusterPort := strconv.Itoa(int(Ports.Port))
			serverPort := strconv.Itoa(int(Ports.TargetPort))
			for i, serverIP := range oldServers {
				probility := strconv.FormatFloat(float64(1/(oldNum-i)), 'f', 4, 64)
				fmt.Printf("Old Probility: %v\n", probility)
				im.DeleteRules(serviceArg.Spec.ClusterIP, serverIP, clusterPort, serverPort, probility)
			}
			for i, serverIP := range newServers {
				probility := strconv.FormatFloat(float64(1/(newNum-i)), 'f', 4, 64)
				fmt.Printf("New Probility: %v\n", probility)
				im.AddRules(serviceArg.Spec.ClusterIP, serverIP, clusterPort, serverPort, probility)
			}
		}
	}

}
