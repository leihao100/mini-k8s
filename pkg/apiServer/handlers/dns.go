package handlers

import (
	"MiniK8S/pkg/api/types"

	"github.com/gin-gonic/gin"
)

func HandleGetDNSs(c *gin.Context) {
	HandleGetApiObjects(c, types.DeploymentObjectType)
}
func HandleGetDNS(c *gin.Context) {
	HandleGetApiObject(c, types.DeploymentObjectType)
}
func HandleCreateDNS(c *gin.Context) {
	HandleCreateApiObject(c, types.DeploymentObjectType)
}
func HandleDeleteDNSs(c *gin.Context) {
	HandleDeleteApiObjects(c, types.DeploymentObjectType)
}
func HandleDeleteDNS(c *gin.Context) {
	HandleDeleteApiObject(c, types.DeploymentObjectType)
}
func HandleWatchDNSs(c *gin.Context) {
	HandleWatchApiObjects(c, types.DeploymentObjectType)
}
func HandleWatchDNS(c *gin.Context) {
	HandleWatchApiObject(c, types.DeploymentObjectType)
}
func HandleGetDNSStatus(c *gin.Context) {
	HandleGetApiObjectStatus(c, types.DeploymentObjectType)
}
func HandleModifyDNSStatus(c *gin.Context) {
	HandleModifyApiObjectStatus(c, types.DeploymentObjectType)
}
