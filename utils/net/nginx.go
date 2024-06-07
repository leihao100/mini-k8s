package net

import (
	"MiniK8S/pkg/api/config"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const DefaultNginxConfigPath = "/usr/local/nginx/conf/nginx.conf"
const NginxConfigPath = DefaultNginxConfigPath

const nginxConfTemplate = `
#minik8s-{{.Metadata.Uid}}
#worker_processes 1;
#events { worker_connections 1024; }
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
#minik8s-end
`

//func CreateNginxConfig(dns config.DNS) error {
//	// 创建一个配置文件
//	configFile, err := os.Create(NginxConfigPath)
//	if err != nil {
//		return fmt.Errorf("error creating config file: %v", err)
//	}
//	defer configFile.Close()
//
//	// 使用模板生成配置文件内容
//	tmpl, err := template.New("nginxConf").Parse(nginxConfTemplate)
//	if err != nil {
//		return fmt.Errorf("error parsing template: %v", err)
//	}
//
//	err = tmpl.Execute(configFile, dns)
//	if err != nil {
//		return fmt.Errorf("error executing template: %v", err)
//	}
//
//	return nil
//}

func GenerateNginxConfig(dns config.DNS) {
	file, err := os.Open(NginxConfigPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	stringbuilder := strings.Builder{}
	tmpl, err := template.New("nginxConf").Parse(nginxConfTemplate)

	err = tmpl.Execute(&stringbuilder, dns)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "#add here" {
			lines = append(lines, line)
			lines = append(lines, stringbuilder.String())
			continue
		}
		lines = append(lines, line)
	}

	file, err = os.Create(NginxConfigPath)

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			//return fmt.Errorf("error writing to file: %v", err)
		}
	}

	RunNginx()
	fmt.Println("Additional content appended to file successfully")
}

func RemoveNginxConfig(dns config.DNS) error {
	start := "#minik8s-" + dns.Metadata.Uid.String()
	end := "#minik8s-end"
	file, err := os.Open(NginxConfigPath)
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
	file, err = os.Create(NginxConfigPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}
	RunNginx()

	return nil

}

func RunNginx() {
	// 测试Nginx配置
	cmd := exec.Command("nginx", "-t")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error testing Nginx config:", err)
		return
	}

	// 启动Nginx
	cmd = exec.Command("nginx", "-s", "reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error starting Nginx:", err)
		return
	}

}
