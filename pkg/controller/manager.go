package controller

import (
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/controller/cache"
	"MiniK8S/pkg/controller/deployment"
	"MiniK8S/pkg/controller/hpa"
	"MiniK8S/pkg/controller/pod"
	"MiniK8S/pkg/controller/storage"
	"MiniK8S/pkg/controller/prometheus"
	"context"
	"fmt"
)

type ControllerManager struct {
	podInformer                   *cache.Informer
	nodeInformer                  *cache.Informer
	deploymentInformer            *cache.Informer
	serviceInformer               *cache.Informer
	hpaInformer                   *cache.Informer
	dnsInformer                   *cache.Informer
	persistentVolumeInformer      *cache.Informer
	persistentVolumeClaimInformer *cache.Informer
	storageClassInformer          *cache.Informer

	podClient                   *apiClient.Client
	nodeClient                  *apiClient.Client
	deploymentClient            *apiClient.Client
	serviceClient               *apiClient.Client
	hpaClient                   *apiClient.Client
	dnsClient                   *apiClient.Client
	persistentVolumeClient      *apiClient.Client
	persistentVolumeClaimClient *apiClient.Client
	storageClassClient          *apiClient.Client

	podController        *pod.PodController
	deploymentController *deployment.DeploymentController
	hpaController        *hpa.HpaController
	storageController    *storage.StorageController
	proController        *prometheus.PrometheusController
}

func NewControllerManager() *ControllerManager {
	fmt.Println("[controller] NewControllerManager")
	podcli, podinf := cache.NewDefaultInformerAndCli(types.PodObjectType)
	nodecli, nodeinf := cache.NewDefaultInformerAndCli(types.NodeObjectType)
	deploymentcli, deploymentinf := cache.NewDefaultInformerAndCli(types.DeploymentObjectType)
	servicecli, serviceinf := cache.NewDefaultInformerAndCli(types.ServiceObjectType)
	hpacli, hpainf := cache.NewDefaultInformerAndCli(types.HorizontalPodAutoscalerObjectType)
	dnscli, dnsinf := cache.NewDefaultInformerAndCli(types.DnsObjectType)
	pvcli, pvinf := cache.NewDefaultInformerAndCli(types.PersistentVolumeObjectType)
	pvccli, pvcinf := cache.NewDefaultInformerAndCli(types.PersistentVolumeClaimObjectType)
	sccli, scinf := cache.NewDefaultInformerAndCli(types.StorageClassObjectType)

	dpController := deployment.NewController(podinf, deploymentinf, podcli, dnscli)
	hpaController := hpa.NewController(podinf, hpainf, deploymentinf, podcli, hpacli, deploymentcli, nodecli)
	storageController := storage.NewController(scinf, pvinf, pvcinf, sccli, pvcli, pvccli)
	prController := prometheus.NewPrometheusController(podcli, nodecli, podinf, nodeinf)
	return &ControllerManager{
		podClient:                   podcli,
		nodeClient:                  nodecli,
		deploymentClient:            deploymentcli,
		serviceClient:               servicecli,
		hpaClient:                   hpacli,
		dnsClient:                   dnscli,
		storageClassClient:          sccli,
		persistentVolumeClient:      pvcli,
		persistentVolumeClaimClient: pvccli,

		podInformer:                   podinf,
		nodeInformer:                  nodeinf,
		deploymentInformer:            deploymentinf,
		serviceInformer:               serviceinf,
		hpaInformer:                   hpainf,
		dnsInformer:                   dnsinf,
		storageClassInformer:          scinf,
		persistentVolumeInformer:      pvinf,
		persistentVolumeClaimInformer: pvcinf,

		deploymentController: dpController,
		hpaController:        hpaController,
		storageController:    storageController,
		proController:        prController,
	}
}

func (cm *ControllerManager) Run(ctx context.Context, cancel context.CancelFunc) {
	fmt.Println("[controller] Starting controller")
	stopCh := make(chan struct{})
	cm.podInformer.Run(stopCh)
	cm.nodeInformer.Run(stopCh)
	cm.deploymentInformer.Run(stopCh)
	cm.serviceInformer.Run(stopCh)
	cm.hpaInformer.Run(stopCh)
	cm.dnsInformer.Run(stopCh)
	cm.persistentVolumeInformer.Run(stopCh)
	cm.persistentVolumeClaimInformer.Run(stopCh)
	cm.storageClassInformer.Run(stopCh)

	cm.deploymentController.Run(ctx, cancel)
	cm.hpaController.Run(ctx, cancel)
	cm.storageController.Run(ctx, cancel)
}
