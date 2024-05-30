package handlers

import (
	"MiniK8S/pkg/api/types"
	"github.com/gin-gonic/gin"
)

func HandleGetHeartbeats(c *gin.Context) {
	HandleGetApiObjects(c, types.HeartbeatObjectType)
}

func HandleGetHeartbeat(c *gin.Context) {
	HandleGetApiObject(c, types.HeartbeatObjectType)
}

func HandleCreateHeartbeat(c *gin.Context) {
	HandleCreateApiObject(c, types.HeartbeatObjectType)
}

func HandleModifyHeartbeat(c *gin.Context) {
	HandleModifyApiObject(c, types.HeartbeatObjectType)
}

func HandleDeleteHeartbeats(c *gin.Context) {
	HandleDeleteApiObjects(c, types.HeartbeatObjectType)
}

func HandleDeleteHeartbeat(c *gin.Context) {
	HandleDeleteApiObject(c, types.HeartbeatObjectType)
}

func HandleWatchHeartbeats(c *gin.Context) {
	HandleWatchApiObjects(c, types.HeartbeatObjectType)
}

func HandleWatchHeartbeat(c *gin.Context) {
	HandleWatchApiObject(c, types.HeartbeatObjectType)
}

func HandleGetHeartbeatStatus(c *gin.Context) {
	HandleGetApiObjectStatus(c, types.HeartbeatObjectType)
}

func HandleModifyHeartbeatStatus(c *gin.Context) {
	HandleModifyApiObjectStatus(c, types.HeartbeatObjectType)
}
