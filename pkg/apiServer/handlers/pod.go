package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleGetPods(c *gin.Context) {
	c.JSON(http.StatusOK, "I'm OK")
}
func HandleGetPod(c *gin.Context) {

}
func HandleCreatePod(c *gin.Context) {

}
func HandleModifyPod(c *gin.Context) {

}
func HandleDeletePods(c *gin.Context) {

}
func HandleDeletePod(c *gin.Context) {

}
func HandleWatchPods(c *gin.Context) {

}
func HandleWatchPod(c *gin.Context) {

}
func HandleGetPodStatus(c *gin.Context) {

}
func HandleModifyPodStatus(c *gin.Context) {

}
