package kubectl

import (
	"MiniK8S/pkg/api/types"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func parseResourceType(ty string) (types.ApiObjectType, error) {
	ty = strings.ToLower(ty)
	switch ty {
	case "pod", "pods":
		return types.PodObjectType, nil
	case "node", "nodes", "no":
		return types.NodeObjectType, nil
	case "service", "svc":
		return types.ServiceObjectType, nil
	case "deployment", "deployments", "deploy":
		return types.DeploymentObjectType, nil
	case "hpa", "hpas":
		return types.HorizontalPodAutoscalerObjectType, nil
	case "storageclass", "storageclasses":
		return types.StorageClassObjectType, nil
	case "persistentvolume", "persistentvolumes":
		return types.PersistentVolumeObjectType, nil
	case "persistentvolumeclaim", "persistentvolumeclaims":
		return types.PersistentVolumeClaimObjectType, nil
	case "dns", "dnss":
		return types.DnsObjectType, nil
	case "job", "jobs":
		return types.JobObjectType, nil

	default:
		errMsg := fmt.Sprintf("No apiObjectType name %s", ty)
		return types.ErrorObjectType, errors.New(errMsg)
	}
}

func getFileName() (string, error) {
	filename, err := rootCmd.PersistentFlags().GetString("filename")
	if err != nil {
		return "", err
	}
	return filename, nil
}

func readFile(file string) (jsonData []byte, ty string, err error) {
	yamlData, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("[kubectl] Error: %v", err)
		return nil, "", err
	}

	var data map[string]interface{}
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		fmt.Printf("[kubectl] Error: %v", err)
		return nil, "", err
	}

	ty, ok := data["kind"].(string)
	if !ok {
		fmt.Printf("[kubectl] Error: kind field not found or is not a string\n")
		return nil, "", fmt.Errorf("kind field not found or is not a string")
	}

	jsonData, err = json.Marshal(data)
	if err != nil {
		fmt.Printf("[kubectl] Error: %v", err)
		return nil, "", err
	}

	return jsonData, ty, nil
}
