package net

import (
	"fmt"
	"os"
	"text/template"
)

const DefaultPrometheusPath = "prometheus.yml"
const prometheusTemplate = `
#name-{{.Name}}
  - job_name: '{{.Name}}'

    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s

    static_configs:
      - targets: ['{{.IP}}']
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
