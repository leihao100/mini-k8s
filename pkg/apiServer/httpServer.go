package apiserver

import (
	"MiniK8S/pkg/apiServer/handlers"

	"github.com/gin-gonic/gin"
)

type HttpServer interface {
	Run(addr string) (err error)
	BindHandlers()
}

func NewHttpServer() HttpServer {
	return &httpServer{
		router: gin.Default(),
	}
}

type httpServer struct {
	router *gin.Engine
}

func (h httpServer) Run(addr string) (err error) {
	return h.router.Run(addr)
}

func (h httpServer) BindHandlers() {
	api := h.router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			clear := v1.Group("/clear")
			{
				clear.DELETE("/", handlers.HandleClear) // DELETE /api/v1/clear/
			}

			pods := v1.Group("/pods")
			{
				pods.GET("/get", handlers.HandleGetPods)                  // GET /api/v1/pods/get
				pods.GET("/get/:name", handlers.HandleGetPod)             // GET /api/v1/pods/get/:name
				pods.PUT("/create", handlers.HandleCreatePod)             // POST /api/v1/pods/create
				pods.PUT("/create/:name", handlers.HandleCreatePod)       // PUT /api/v1/pods/create/:name
				pods.DELETE("/delete", handlers.HandleDeletePods)         // DELETE /api/v1/pods/delete
				pods.DELETE("/delete/:name", handlers.HandleDeletePod)    // DELETE /api/v1/pods/delete/:name
				pods.GET("/watch", handlers.HandleWatchPods)              // GET /api/v1/pods/watch
				pods.GET("/watch/:name", handlers.HandleWatchPod)         // GET /api/v1/pods/watch/:name
				pods.GET("/status/:name", handlers.HandleGetPodStatus)    // GET /api/v1/pods/status/:name
				pods.PUT("/status/:name", handlers.HandleModifyPodStatus) // PUT /api/v1/pods/status/:name
			}
			services := v1.Group("/services")
			{
				services.GET("/get", handlers.HandleGetServices)                  // GET /api/v1/services/get
				services.GET("/get/:name", handlers.HandleGetService)             // GET /api/v1/services/get/:name
				services.PUT("/create", handlers.HandleCreateService)             // POST /api/v1/services/create
				services.PUT("/create/:name", handlers.HandleCreateService)       // PUT /api/v1/services/create/:name
				services.DELETE("/delete", handlers.HandleDeleteServices)         // DELETE /api/v1/services/delete
				services.DELETE("/delete/:name", handlers.HandleDeleteService)    // DELETE /api/v1/services/delete/:name
				services.GET("/watch", handlers.HandleWatchServices)              // GET /api/v1/services/watch
				services.GET("/watch/:name", handlers.HandleWatchService)         // GET /api/v1/services/watch/:name
				services.GET("/status/:name", handlers.HandleGetServiceStatus)    // GET /api/v1/services/status/:name
				services.PUT("/status/:name", handlers.HandleModifyServiceStatus) // PUT /api/v1/services/status/:name
			}
			deployments := v1.Group("/deployments")
			{
				deployments.GET("/get", handlers.HandleGetDeployments)                  // GET /api/v1/deployments/get
				deployments.GET("/get/:name", handlers.HandleDeleteDeployment)          // GET /api/v1/deployments/get/:name
				deployments.PUT("/create", handlers.HandleCreateDeployment)             // POST /api/v1/deployments/create
				deployments.PUT("/create/:name", handlers.HandleCreateDeployment)       // PUT /api/v1/deployments/create/:name
				deployments.DELETE("/delete", handlers.HandleDeleteDeployments)         // DELETE /api/v1/deployments/delete
				deployments.DELETE("/delete/:name", handlers.HandleDeleteDeployment)    // DELETE /api/v1/deployments/delete/:name
				deployments.GET("/watch", handlers.HandleWatchDeployments)              // GET /api/v1/deployments/watch
				deployments.GET("/watch/:name", handlers.HandleWatchDeployment)         // GET /api/v1/deployments/watch/:name
				deployments.GET("/status/:name", handlers.HandleGetDeploymentStatus)    // GET /api/v1/deployments/status/:name
				deployments.PUT("/status/:name", handlers.HandleModifyDeploymentStatus) // PUT /api/v1/deployments/status/:name
			}
			hpas := v1.Group("/hpas")
			{
				hpas.GET("/get", handlers.HandleGetHPAs)                  // GET /api/v1/hpa/get
				hpas.GET("/get/:name", handlers.HandleGetHPA)             // GET /api/v1/hpa/get/:name
				hpas.PUT("/create", handlers.HandleCreateHPA)             // POST /api/v1/hpa/create
				hpas.PUT("/create/:name", handlers.HandleCreateHPA)       // PUT /api/v1/hpa/create/:name
				hpas.DELETE("/delete", handlers.HandleDeleteHPAs)         // DELETE /api/v1/hpa/delete
				hpas.DELETE("/delete/:name", handlers.HandleDeleteHPA)    // DELETE /api/v1/hpa/delete/:name
				hpas.GET("/watch", handlers.HandleWatchHPAs)              // GET /api/v1/hpa/watch
				hpas.GET("/watch/:name", handlers.HandleWatchHPA)         // GET /api/v1/hpa/watch/:name
				hpas.GET("/status/:name", handlers.HandleGetHPAStatus)    // GET /api/v1/hpa/status/:name
				hpas.PUT("/status/:name", handlers.HandleModifyHPAStatus) // PUT /api/v1/hpa/status/:name
			}
			nodes := v1.Group("/nodes")
			{
				nodes.GET("/get", handlers.HandleGetNodes)                  // GET /api/v1/nodes/get
				nodes.GET("/get/:name", handlers.HandleGetNode)             // GET /api/v1/nodes/get/:name
				nodes.PUT("/create", handlers.HandleCreateNode)             // POST /api/v1/nodes/create
				nodes.PUT("/create/:name", handlers.HandleCreateNode)       // PUT /api/v1/nodes/create/:name
				nodes.DELETE("/delete", handlers.HandleDeleteNodes)         // DELETE /api/v1/nodes/delete
				nodes.DELETE("/delete/:name", handlers.HandleDeleteNode)    // DELETE /api/v1/nodes/delete/:name
				nodes.GET("/watch", handlers.HandleWatchNodes)              // GET /api/v1/nodes/watch
				nodes.GET("/watch/:name", handlers.HandleWatchNode)         // GET /api/v1/nodes/watch/:name
				nodes.GET("/status/:name", handlers.HandleGetNodeStatus)    // GET /api/v1/nodes/status/:name
				nodes.PUT("/status/:name", handlers.HandleModifyNodeStatus) // PUT /api/v1/nodes/status/:name
			}
			heartbeats := v1.Group("/heartbeats")
			{
				heartbeats.GET("/get", handlers.HandleGetHeartbeats)                  // GET /api/v1/heartbeats/get
				heartbeats.GET("/get/:name", handlers.HandleGetHeartbeat)             // GET /api/v1/heartbeats/get/:name
				heartbeats.PUT("/create", handlers.HandleCreateHeartbeat)             // POST /api/v1/heartbeats/create
				heartbeats.PUT("/create/:name", handlers.HandleCreateHeartbeat)       // PUT /api/v1/heartbeats/create/:name
				heartbeats.DELETE("/delete", handlers.HandleDeleteHeartbeats)         // DELETE /api/v1/heartbeats/delete
				heartbeats.DELETE("/delete/:name", handlers.HandleDeleteHeartbeat)    // DELETE /api/v1/heartbeats/delete/:name
				heartbeats.GET("/watch", handlers.HandleWatchHeartbeats)              // GET /api/v1/heartbeats/watch
				heartbeats.GET("/watch/:name", handlers.HandleWatchHeartbeat)         // GET /api/v1/heartbeats/watch/:name
				heartbeats.GET("/status/:name", handlers.HandleGetHeartbeatStatus)    // GET /api/v1/heartbeats/status/:name
				heartbeats.PUT("/status/:name", handlers.HandleModifyHeartbeatStatus) // PUT /api/v1/heartbeats/status/:name
			}
			dns := v1.Group("/dnss")
			{
				dns.GET("/get", handlers.HandleGetDNSs)                  // GET /api/v1/heartbeats/get
				dns.GET("/get/:name", handlers.HandleGetDNS)             // GET /api/v1/heartbeats/get/:name
				dns.PUT("/create", handlers.HandleCreateDNS)             // POST /api/v1/heartbeats/create
				dns.PUT("/create/:name", handlers.HandleCreateDNS)       // PUT /api/v1/heartbeats/create/:name
				dns.DELETE("/delete", handlers.HandleDeleteDNSs)         // DELETE /api/v1/heartbeats/delete
				dns.DELETE("/delete/:name", handlers.HandleDeleteDNS)    // DELETE /api/v1/heartbeats/delete/:name
				dns.GET("/watch", handlers.HandleWatchDNSs)              // GET /api/v1/heartbeats/watch
				dns.GET("/watch/:name", handlers.HandleWatchDNS)         // GET /api/v1/heartbeats/watch/:name
				dns.GET("/status/:name", handlers.HandleGetDNSStatus)    // GET /api/v1/heartbeats/status/:name
				dns.PUT("/status/:name", handlers.HandleModifyDNSStatus) // PUT /api/v1/heartbeats/status/:name
			}
		}
	}
}
