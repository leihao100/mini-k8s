package kubectl

import (
	"MiniK8S/pkg/api/types"
	"errors"
	"fmt"
)

func parseResourceType(ty string) (types.ApiObjectType, error) {
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
	default:
		errMsg := fmt.Sprintf("No ObjectType name %s", ty)
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
