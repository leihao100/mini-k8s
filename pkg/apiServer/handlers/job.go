package handlers

import (
	"MiniK8S/pkg/api/types"

	"github.com/gin-gonic/gin"
)

func HandleGetJobs(c *gin.Context) {
	HandleGetApiObjects(c, types.JobObjectType)
}

func HandleGetJob(c *gin.Context) {
	HandleGetApiObject(c, types.JobObjectType)
}

func HandleCreateJob(c *gin.Context) {
	HandleCreateApiObject(c, types.JobObjectType)
}

func HandleDeleteJobs(c *gin.Context) {
	HandleDeleteApiObjects(c, types.JobObjectType)
}

func HandleDeleteJob(c *gin.Context) {
	HandleDeleteApiObject(c, types.JobObjectType)
}

func HandleWatchJobs(c *gin.Context) {
	HandleWatchApiObjects(c, types.JobObjectType)
}

func HandleWatchJob(c *gin.Context) {
	HandleWatchApiObject(c, types.JobObjectType)
}

func HandleGetJobStatus(c *gin.Context) {
	HandleGetApiObjectStatus(c, types.JobObjectType)
}

func HandleModifyJobStatus(c *gin.Context) {
	HandleModifyApiObjectStatus(c, types.JobObjectType)
}
