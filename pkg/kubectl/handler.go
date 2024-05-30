package kubectl

import (
	"MiniK8S/config"
	apiConfig "MiniK8S/pkg/api/config"
	"MiniK8S/pkg/apiClient"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

func create(cmd *cobra.Command, args []string) {
	filename, err := getFileName()
	if err != nil {
		fmt.Printf("[kubectl] Error: %v\n", err)
		return
	}
	jsonData, resourceType, err := readFile(filename)
	if err != nil {
		fmt.Printf("[kubectl] Error: %v\n", err)
		return
	}
	apiObjectType, err := parseResourceType(resourceType)
	if err != nil {
		fmt.Printf("[kubectl] Error: %v\n", err)
		return
	}
	cli := apiClient.NewRESTClient(apiObjectType)
	apiObject := apiConfig.NewApiObject(apiObjectType)
	err = apiObject.JsonUnmarshal(jsonData)
	if err != nil {
		fmt.Printf("[kubectl] Error: %v\n", err)
	}
	serverPath := cli.BuildURL(apiClient.Create)
	cli.Post(serverPath, jsonData)
}

func delete(cmd *cobra.Command, args []string) {
	resourceType := args[0]
	name := args[1]

	apiObjectType, err := parseResourceType(resourceType)
	if err != nil {
		fmt.Printf("[kubectl] Error: %v\n", err)
		return
	}
	cli := apiClient.NewRESTClient(apiObjectType)
	serverPath := cli.BuildURL(apiClient.Delete) + name
	cli.Delete(serverPath, nil)
}

func get(cmd *cobra.Command, args []string) {
	resourceType := args[0]
	apiObjectType, err := parseResourceType(resourceType)
	if err != nil {
		fmt.Printf("[kubectl] Error: %v\n", err)
		return
	}
	cli := apiClient.NewRESTClient(apiObjectType)
	if len(args) == 1 {
		serverPath := cli.BuildURL(apiClient.Get)
		res := cli.Get(serverPath, nil)
		var jsonData []byte
		res.Read(jsonData)
		apiObjectList := apiConfig.NewApiObjectList(apiObjectType)
		err := apiObjectList.JsonUnmarshal(jsonData)
		if err != nil {
			fmt.Printf("[kubectl] Error: %v\n", err)
		}
		apiObjectList.Info()
	} else if len(args) >= 2 {
		name := args[1]
		serverPath := cli.BuildURL(apiClient.Get) + name
		res := cli.Get(serverPath, nil)
		var jsonData []byte
		res.Read(jsonData)
		apiObject := apiConfig.NewApiObject(apiObjectType)
		err := apiObject.JsonUnmarshal(jsonData)
		if err != nil {
			fmt.Printf("[kubectl] Error: %v\n", err)
		}
		apiObject.Info()
	}
}

func describe(cmd *cobra.Command, args []string) {
	resourceType := args[0]
	apiObjectType, err := parseResourceType(resourceType)
	if err != nil {
		fmt.Printf("[kubectl] Error: %v\n", err)
	}
	cli := apiClient.NewRESTClient(apiObjectType)
	if len(args) == 1 {
		serverPath := cli.BuildURL(apiClient.Get)
		res := cli.Get(serverPath, nil)
		var jsonData []byte
		res.Read(jsonData)
		yamlData, err := yaml.JSONToYAML(jsonData)
		if err != nil {
			fmt.Printf("[kubectl] Error: %v\n", err)
		}
		fmt.Printf("%v\n", yamlData)
	} else if len(args) >= 2 {
		name := args[1]
		serverPath := cli.BuildURL(apiClient.Get) + name
		res := cli.Get(serverPath, nil)
		var jsonData []byte
		res.Read(jsonData)
		yamlData, err := yaml.JSONToYAML(jsonData)
		if err != nil {
			fmt.Printf("[kubectl] Error: %v\n", err)
		}
		fmt.Printf("%v\n", yamlData)
	}
}

func clear(cmd *cobra.Command, args []string) {
	serverPath := config.ApiServerHost() + config.ApiServerPort() + "/api/v1/clear"
	req, _ := http.NewRequest("DELETE", serverPath, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("[kubectl] Error: %v\n", err)
	}
	str, _ := io.ReadAll(res.Body)

	var result map[string]interface{}
	err = json.Unmarshal(str, &result)
	if err != nil {
		fmt.Printf("[kubectl] Error parsing JSON: %v\n", err)
		return
	}

	status, ok := result["status"].(string)
	if !ok {
		fmt.Printf("[kubectl] Error: status field not found or is not a string\n")
		return
	}
	if status == "ERR" {
		fmt.Printf("[kubectl] clear failed!\n")

	}
	fmt.Printf("[kubectl] clear successfully\n")
}
