package handlers

import (
	"MiniK8S/pkg/api/types"

	"github.com/gin-gonic/gin"
)

func HandleGetDNSs(c *gin.Context) {
	HandleGetApiObjects(c, types.DnsObjectType)
}

func HandleGetDNS(c *gin.Context) {
	HandleGetApiObject(c, types.DnsObjectType)
}

func HandleCreateDNS(c *gin.Context) {
	HandleCreateApiObject(c, types.DnsObjectType)
}

func HandleDeleteDNSs(c *gin.Context) {
	HandleDeleteApiObjects(c, types.DnsObjectType)
}

func HandleDeleteDNS(c *gin.Context) {
	HandleDeleteApiObject(c, types.DnsObjectType)
}

func HandleWatchDNSs(c *gin.Context) {
	HandleWatchApiObjects(c, types.DnsObjectType)
}

func HandleWatchDNS(c *gin.Context) {
	HandleWatchApiObject(c, types.DnsObjectType)
}

func HandleGetDNSStatus(c *gin.Context) {
	HandleGetApiObjectStatus(c, types.DnsObjectType)
}

func HandleModifyDNSStatus(c *gin.Context) {
	HandleModifyApiObjectStatus(c, types.DnsObjectType)
}
