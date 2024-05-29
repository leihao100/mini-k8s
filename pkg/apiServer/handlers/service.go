package handlers

import (
	"MiniK8S/pkg/api/types"

	"github.com/gin-gonic/gin"
)

func HandleGetServices(c *gin.Context) {
	HandleGetApiObjects(c, types.ServiceObjectType)
}

func HandleGetService(c *gin.Context) {
	HandleGetApiObject(c, types.ServiceObjectType)
}

func HandleCreateService(c *gin.Context) {
	HandleCreateApiObject(c, types.ServiceObjectType)
}

func HandleModifyService(c *gin.Context) {
	HandleModifyApiObject(c, types.ServiceObjectType)
}

func HandleDeleteServices(c *gin.Context) {
	HandleDeleteApiObjects(c, types.ServiceObjectType)
}

func HandleDeleteService(c *gin.Context) {
	HandleDeleteApiObject(c, types.ServiceObjectType)
}

func HandleWatchServices(c *gin.Context) {
	HandleWatchApiObjects(c, types.ServiceObjectType)
}

func HandleWatchService(c *gin.Context) {
	HandleWatchApiObject(c, types.ServiceObjectType)
}

func HandleGetServiceStatus(c *gin.Context) {
	HandleGetApiObjectStatus(c, types.ServiceObjectType)
}

func HandleModifyServiceStatus(c *gin.Context) {
	HandleModifyApiObjectStatus(c, types.ServiceObjectType)
}
