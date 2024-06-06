package storage

import (
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/controller/cache"
)

type StorageController struct {
	storageClassClient          *apiClient.Client
	persistentVolumeClient      *apiClient.Client
	persistentVolumeClaimClient *apiClient.Client
	storageClassInformer        *cache.Informer
}
