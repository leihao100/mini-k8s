package controller

import (
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/controller/cache"
	"MiniK8S/pkg/controller/deployment"
	"MiniK8S/pkg/controller/hpa"
	"MiniK8S/pkg/controller/pod"
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
	podcli, podinf := cache.NewDefaultInformerAndCli(types.PodObjectType)
	nodecli, nodeinf := cache.NewDefaultInformerAndCli(types.NodeObjectType)
	deploymentcli, deploymentinf := cache.NewDefaultInformerAndCli(types.DeploymentObjectType)
	servicecli, serviceinf := cache.NewDefaultInformerAndCli(types.ServiceObjectType)
	hpacli, hpainf := cache.NewDefaultInformerAndCli(types.HorizontalPodAutoscalerObjectType)
	dnscli, dnsinf := cache.NewDefaultInformerAndCli(types.DnsObjectType)
	dpController := deployment.NewController(podinf, deploymentinf, podcli, dnscli)
	//hpaController := hpa.NewController(podinf,hpainf)
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
	}
}
