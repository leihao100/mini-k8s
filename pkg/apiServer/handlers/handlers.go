package handlers

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/types"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetApiObjects(c *gin.Context, ty types.ApiObjectType) {
	switch ty {
	case types.PodObjectType:
	case types.DeploymentObjectType:
	}
}

func HandleGetApiObject(c *gin.Context, ty types.ApiObjectType) {
	switch ty {
	case types.PodObjectType:
	case types.DeploymentObjectType:
	}
}

func HandleCreateApiObject(c *gin.Context, ty types.ApiObjectType) {
	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
	}
	newApiObject := config.NewApiObject(ty)
	if newApiObject {

	}
	err = newApiObject.JsonUnmarshal(buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
}

func HandleModifyApiObject(c *gin.Context, ty types.ApiObjectType) {
	switch ty {
	case types.PodObjectType:
	case types.DeploymentObjectType:
	}
}

func HandleDeleteApiObjects(c *gin.Context, ty types.ApiObjectType) {
	switch ty {
	case types.PodObjectType:
	case types.DeploymentObjectType:
	}
}

func HandleDeleteApiObject(c *gin.Context, ty types.ApiObjectType) {
	switch ty {
	case types.PodObjectType:
	case types.DeploymentObjectType:
	}
}

func HandleWatchApiObjects(c *gin.Context, ty types.ApiObjectType) {
	switch ty {
	case types.PodObjectType:
	case types.DeploymentObjectType:
	}
}

func HandleWatchApiObject(c *gin.Context, ty types.ApiObjectType) {
	switch ty {
	case types.PodObjectType:
	case types.DeploymentObjectType:
	}
}

func HandleGetApiObjectStatus(c *gin.Context, ty types.ApiObjectType) {
	switch ty {
	case types.PodObjectType:
	case types.DeploymentObjectType:
	}
}

func HandleModifyApiObjectStatus(c *gin.Context, ty types.ApiObjectType) {
	switch ty {
	case types.PodObjectType:
	case types.DeploymentObjectType:
	}
}
