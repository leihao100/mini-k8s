package handlers

import (
	"MiniK8S/pkg/api/types"

	"github.com/gin-gonic/gin"
)

func HandleGetHPAs(c *gin.Context) {
	HandleGetApiObjects(c, types.HorizontalPodAutoscalerObjectType)
}

func HandleGetHPA(c *gin.Context) {
	HandleGetApiObject(c, types.HorizontalPodAutoscalerObjectType)
}

func HandleCreateHPA(c *gin.Context) {
	HandleCreateApiObject(c, types.HorizontalPodAutoscalerObjectType)
}

func HandleDeleteHPAs(c *gin.Context) {
	HandleDeleteApiObjects(c, types.HorizontalPodAutoscalerObjectType)
}

func HandleDeleteHPA(c *gin.Context) {
	HandleDeleteApiObject(c, types.HorizontalPodAutoscalerObjectType)
}

func HandleWatchHPAs(c *gin.Context) {
	HandleWatchApiObjects(c, types.HorizontalPodAutoscalerObjectType)
}

func HandleWatchHPA(c *gin.Context) {
	HandleWatchApiObject(c, types.HorizontalPodAutoscalerObjectType)
}

func HandleGetHPAStatus(c *gin.Context) {
	HandleGetApiObjectStatus(c, types.HorizontalPodAutoscalerObjectType)
}

func HandleModifyHPAStatus(c *gin.Context) {
	HandleModifyApiObjectStatus(c, types.HorizontalPodAutoscalerObjectType)
}
