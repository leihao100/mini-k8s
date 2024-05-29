package handlers

import (
	"MiniK8S/pkg/api/types"
	"github.com/gin-gonic/gin"
)

func HandleCreateHeartbeat(c *gin.Context) {
	HandleGetApiObjects(c, types.HeartbeatObjectType)
}
