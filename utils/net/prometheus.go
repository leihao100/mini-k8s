package net

import (
	"bufio"
	"fmt"
	"os"
	"text/template"
)

const DefaultPrometheusPath = "prometheus.yml"
const prometheusTemplate = `
#minik8s-{{.Name}}
- targets:
  - '{{.IP}}'
  labels:
    job: '{{.Name}}'
#minik8s-end
`

type pro struct {
	Name string
	IP   string
}

func AddPrometheus(name, ip string) {
	file, err := os.OpenFile(DefaultPrometheusPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	conf := pro{
		Name: name,
		IP:   ip,
	}
	tmpl, err := template.New("prometheusConfig").Parse(prometheusTemplate)
	err = tmpl.Execute(file, conf)
}

func RemovePrometheus(name string) error {
	start := "#minik8s-" + name
	end := "#minik8s-end"
	file, err := os.Open(DefaultPrometheusPath)
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
	file, err = os.Create(DefaultPrometheusPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	return nil

}
