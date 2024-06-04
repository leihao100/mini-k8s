package controller

import (
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/controller/cache"
	"MiniK8S/pkg/controller/deployment"
	"MiniK8S/pkg/controller/hpa"
	"MiniK8S/pkg/controller/pod"
	"context"
	"fmt"
)

type ControllerManager struct {
	podInformer        *cache.Informer
	nodeInformer       *cache.Informer
	deploymentInformer *cache.Informer
	serviceInformer    *cache.Informer
	hpaInformer        *cache.Informer
	dnsInformer        *cache.Informer

	podClient        *apiClient.Client
	nodeClient       *apiClient.Client
	deploymentClient *apiClient.Client
	serviceClient    *apiClient.Client
	hpaClient        *apiClient.Client
	dnsClient        *apiClient.Client

	podController        *pod.PodController
	deploymentController *deployment.DeploymentController
	hpaController        *hpa.HpaController
}

func NewControllerManager() *ControllerManager {
	fmt.Println("[controller] NewControllerManager")
	podcli, podinf := cache.NewDefaultInformerAndCli(types.PodObjectType)
	nodecli, nodeinf := cache.NewDefaultInformerAndCli(types.NodeObjectType)
	deploymentcli, deploymentinf := cache.NewDefaultInformerAndCli(types.DeploymentObjectType)
	servicecli, serviceinf := cache.NewDefaultInformerAndCli(types.ServiceObjectType)
	hpacli, hpainf := cache.NewDefaultInformerAndCli(types.HorizontalPodAutoscalerObjectType)
	dnscli, dnsinf := cache.NewDefaultInformerAndCli(types.DnsObjectType)
	dpController := deployment.NewController(podinf, deploymentinf, podcli, dnscli)
	hpaController := hpa.NewController(podinf, hpainf, podcli, hpacli, deploymentcli)
	return &ControllerManager{
		podClient:        podcli,
		nodeClient:       nodecli,
		deploymentClient: deploymentcli,
		serviceClient:    servicecli,
		hpaClient:        hpacli,
		dnsClient:        dnscli,

		podInformer:        podinf,
		nodeInformer:       nodeinf,
		deploymentInformer: deploymentinf,
		serviceInformer:    serviceinf,
		hpaInformer:        hpainf,
		dnsInformer:        dnsinf,

		deploymentController: dpController,
		hpaController:        hpaController,
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

	cm.deploymentController.Run(ctx, cancel)
	cm.hpaController.Run(ctx, cancel)
}
