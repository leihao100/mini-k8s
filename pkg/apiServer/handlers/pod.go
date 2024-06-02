package handlers

import (
	"MiniK8S/pkg/api/types"

	"github.com/gin-gonic/gin"
)

func HandleGetPods(c *gin.Context) {
	HandleGetApiObjects(c, types.PodObjectType)
}

func HandleGetPod(c *gin.Context) {
	HandleGetApiObject(c, types.PodObjectType)
}

func HandleCreatePod(c *gin.Context) {
	HandleCreateApiObject(c, types.PodObjectType)
}

func HandleDeletePods(c *gin.Context) {
	HandleDeleteApiObjects(c, types.PodObjectType)
}

func HandleDeletePod(c *gin.Context) {
	HandleDeleteApiObject(c, types.PodObjectType)
}

func HandleWatchPods(c *gin.Context) {
	HandleWatchApiObjects(c, types.PodObjectType)
}

func HandleWatchPod(c *gin.Context) {
	HandleWatchApiObject(c, types.PodObjectType)
}

func HandleGetPodStatus(c *gin.Context) {
	HandleGetApiObjectStatus(c, types.PodObjectType)
}

func HandleModifyPodStatus(c *gin.Context) {
	HandleModifyApiObjectStatus(c, types.PodObjectType)
}
