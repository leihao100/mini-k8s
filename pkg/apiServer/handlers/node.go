package handlers

import (
	"MiniK8S/pkg/api/types"

	"github.com/gin-gonic/gin"
)

func HandleGetNodes(c *gin.Context) {
	HandleGetApiObjects(c, types.NodeObjectType)
}

func HandleGetNode(c *gin.Context) {
	HandleGetApiObject(c, types.NodeObjectType)
}

func HandleCreateNode(c *gin.Context) {
	HandleCreateApiObject(c, types.NodeObjectType)
}

func HandleDeleteNodes(c *gin.Context) {
	HandleDeleteApiObjects(c, types.NodeObjectType)
}

func HandleDeleteNode(c *gin.Context) {
	HandleDeleteApiObject(c, types.NodeObjectType)
}

func HandleWatchNodes(c *gin.Context) {
	HandleWatchApiObjects(c, types.NodeObjectType)
}

func HandleWatchNode(c *gin.Context) {
	HandleWatchApiObject(c, types.NodeObjectType)
}

func HandleGetNodeStatus(c *gin.Context) {
	HandleGetApiObjectStatus(c, types.NodeObjectType)
}

func HandleModifyNodeStatus(c *gin.Context) {
	HandleModifyApiObjectStatus(c, types.NodeObjectType)
}
