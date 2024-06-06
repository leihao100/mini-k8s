package handlers

import (
	"MiniK8S/pkg/api/types"
	"github.com/gin-gonic/gin"
)

// StorageClass handlers
func HandleGetStorageClasses(c *gin.Context) {
	HandleGetApiObjects(c, types.StorageClassObjectType)
}

func HandleGetStorageClass(c *gin.Context) {
	HandleGetApiObject(c, types.StorageClassObjectType)
}

func HandleCreateStorageClass(c *gin.Context) {
	HandleCreateApiObject(c, types.StorageClassObjectType)
}

func HandleDeleteStorageClass(c *gin.Context) {
	HandleDeleteApiObject(c, types.StorageClassObjectType)
}

func HandleWatchStorageClasses(c *gin.Context) {
	HandleWatchApiObjects(c, types.StorageClassObjectType)
}

func HandleWatchStorageClass(c *gin.Context) {
	HandleWatchApiObject(c, types.StorageClassObjectType)
}

func HandleGetStorageClassStatus(c *gin.Context) {
	HandleGetApiObjectStatus(c, types.StorageClassObjectType)
}

func HandleModifyStorageClassStatus(c *gin.Context) {
	HandleModifyApiObjectStatus(c, types.StorageClassObjectType)
}

// PersistentVolume handlers
func HandleGetPersistentVolumes(c *gin.Context) {
	HandleGetApiObjects(c, types.PersistentVolumeObjectType)
}

func HandleGetPersistentVolume(c *gin.Context) {
	HandleGetApiObject(c, types.PersistentVolumeObjectType)
}

func HandleCreatePersistentVolume(c *gin.Context) {
	HandleCreateApiObject(c, types.PersistentVolumeObjectType)
}

func HandleDeletePersistentVolume(c *gin.Context) {
	HandleDeleteApiObject(c, types.PersistentVolumeObjectType)
}

func HandleWatchPersistentVolumes(c *gin.Context) {
	HandleWatchApiObjects(c, types.PersistentVolumeObjectType)
}

func HandleWatchPersistentVolume(c *gin.Context) {
	HandleWatchApiObject(c, types.PersistentVolumeObjectType)
}

func HandleGetPersistentVolumeStatus(c *gin.Context) {
	HandleGetApiObjectStatus(c, types.PersistentVolumeObjectType)
}

// PersistentVolumeClaim handlers
func HandleGetPersistentVolumeClaims(c *gin.Context) {
	HandleGetApiObjects(c, types.PersistentVolumeClaimObjectType)
}

func HandleGetPersistentVolumeClaim(c *gin.Context) {
	HandleGetApiObject(c, types.PersistentVolumeClaimObjectType)
}

func HandleCreatePersistentVolumeClaim(c *gin.Context) {
	HandleCreateApiObject(c, types.PersistentVolumeClaimObjectType)
}

func HandleDeletePersistentVolumeClaim(c *gin.Context) {
	HandleDeleteApiObject(c, types.PersistentVolumeClaimObjectType)
}

func HandleWatchPersistentVolumeClaims(c *gin.Context) {
	HandleWatchApiObjects(c, types.PersistentVolumeClaimObjectType)
}

func HandleWatchPersistentVolumeClaim(c *gin.Context) {
	HandleWatchApiObject(c, types.PersistentVolumeClaimObjectType)
}

func HandleGetPersistentVolumeClaimStatus(c *gin.Context) {
	HandleGetApiObjectStatus(c, types.PersistentVolumeClaimObjectType)
}

func HandleModifyPersistentVolumeClaimStatus(c *gin.Context) {
	HandleModifyApiObjectStatus(c, types.PersistentVolumeClaimObjectType)
}
