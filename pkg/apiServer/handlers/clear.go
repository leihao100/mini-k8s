package handlers

import (
	"MiniK8S/pkg/etcd"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleClear(c *gin.Context) {
	err := etcd.Clear()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
