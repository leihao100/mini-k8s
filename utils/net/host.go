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
#minik8s-{{.Metadata.Uid}}
192.168.1.10  {{.Spec.HostName}}
#minik8s-end
`

func AddHost(dns *config.DNS) {
	file, err := os.OpenFile(HostConfigPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	tmpl, err := template.New("nginxConf").Parse(HostConfTemplate)

	err = tmpl.Execute(file, dns)

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
	//defer file.Close()
	return nil
}
