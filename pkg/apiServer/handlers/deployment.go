package handlers

import (
	"MiniK8S/pkg/api/types"

	"github.com/gin-gonic/gin"
)

func HandleGetDeployments(c *gin.Context) {
	HandleGetApiObjects(c, types.DeploymentObjectType)
}
func HandleGetDeployment(c *gin.Context) {
	HandleGetApiObject(c, types.DeploymentObjectType)
}
func HandleCreateDeployment(c *gin.Context) {
	HandleCreateApiObject(c, types.DeploymentObjectType)
}
func HandleModifyDeployment(c *gin.Context) {
	HandleModifyApiObject(c, types.DeploymentObjectType)
}
func HandleDeleteDeployments(c *gin.Context) {
	HandleDeleteApiObjects(c, types.DeploymentObjectType)
}
func HandleDeleteDeployment(c *gin.Context) {
	HandleDeleteApiObject(c, types.DeploymentObjectType)
}
func HandleWatchDeployments(c *gin.Context) {
	HandleWatchApiObjects(c, types.DeploymentObjectType)
}
func HandleWatchDeployment(c *gin.Context) {
	HandleWatchApiObject(c, types.DeploymentObjectType)
}
func HandleGetDeploymentStatus(c *gin.Context) {
	HandleGetApiObjectStatus(c, types.DeploymentObjectType)
}
func HandleModifyDeploymentStatus(c *gin.Context) {
	HandleModifyApiObjectStatus(c, types.DeploymentObjectType)
}
