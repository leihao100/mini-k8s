package net

import (
	"MiniK8S/pkg/api/config"
	"bufio"
	"fmt"
	"os"
	"text/template"
)

const HostConfigPath = "/etc/hosts"

const HostConfTemplate = `
#minik8s-{{.Uid}}
{{.Ip}}  {{.HostName}}
#minik8s-end
`

type hostconfig struct {
	Uid      string
	Ip       string
	HostName string
}

func AddHost(dns *config.DNS) {
	file, err := os.OpenFile(HostConfigPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	tmpl, err := template.New("nginxConf").Parse(HostConfTemplate)
	ip, _ := GetLocalIP()
	con := hostconfig{
		Uid:      dns.Metadata.Uid.String(),
		Ip:       ip,
		HostName: dns.Spec.HostName,
	}
	err = tmpl.Execute(file, con)

	fmt.Println("Additional content appended to file successfully")
}
func RemoveHost(dns config.DNS) error {
	start := "#minik8s-" + dns.Metadata.Uid.String()
	end := "#minik8s-end"
	file, err := os.Open(HostConfigPath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()
	flag := false
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if flag == false {
			if line == start {
				flag = true
				continue
			}
		}
		if flag == true {
			if line != end {
				continue
			} else {
				flag = false
				continue
			}
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	file, err = os.Create(HostConfigPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}
	//defer file.Close()
	return nil
}
