package net

import (
	"MiniK8S/pkg/api/config"
	"fmt"
	"os"
	"text/template"
)

const DefaultNginxConfigPath = "/etc/nginx/nginx.conf"
const NginxConfigPath = "nginx.conf"

const nginxConfTemplate = `
#worker_processes 1;

#events { worker_connections 1024; }

http {
    server {
        listen {{.Spec.HostPort}};

        server_name {{.Spec.HostName}};

        {{range .Spec.Path}}
        location {{.ClusterPath}} {
            proxy_pass http://{{.ClusterIP}};
            #proxy_set_header Host $host;
            #proxy_set_header X-Real-IP $remote_addr;
            #proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            #proxy_set_header X-Forwarded-Proto $scheme;
        }
        {{end}}
    }
}
`

func CreateNginxConfig(dns config.DNS) error {
	// 创建一个配置文件
	configFile, err := os.Create(NginxConfigPath)
	if err != nil {
		return fmt.Errorf("error creating config file: %v", err)
	}
	defer configFile.Close()

	// 使用模板生成配置文件内容
	tmpl, err := template.New("nginxConf").Parse(nginxConfTemplate)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	err = tmpl.Execute(configFile, dns)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}

func GenerateNginxConfig(dns config.DNS) {
	file, err := os.OpenFile(NginxConfigPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	tmpl, err := template.New("nginxConf").Parse(nginxConfTemplate)

	err = tmpl.Execute(file, dns)

	fmt.Println("Additional content appended to file successfully")
}
