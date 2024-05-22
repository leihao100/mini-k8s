package apiserver

import "github.com/gin-gonic/gin"

type HttpServer interface {
	Run(addr string) (err error)
	BindHandlers()
}

func New() HttpServer {
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

/*
"/api/pods/"

*/

func (h httpServer) BindHandlers() {
	api := h.router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			clear := v1.Group("/clear")
			{
				clear.DELETE("/")
			}

			pods := v1.Group("/pods")
			{
				pods.GET("/get")             //GET /api/v1/pods/get
				pods.GET("/get/:name")       //GET /api/v1/pods/get/:name
				pods.POST("/create")         //POST /api/v1/pods/get
				pods.PUT("/create/:name")    //PUT /api/v1/pods/get
				pods.DELETE("/delete")       //DELETE /api/v1/pods/delete
				pods.DELETE("/delete/:name") //DELETE /api/v1/pods/get
				pods.GET("/watch")
				pods.GET("/watch/:name")
				pods.GET("/status/:name")
				pods.PUT("/status/:name")
			}
			services := v1.Group("/services")
			{
				services.GET("/get")
				services.GET("/get/:name")
				services.POST("/create")
				services.PUT("/create/:name")
				services.DELETE("/delete")
				services.DELETE("/delete/:name")
				services.GET("/watch")
				services.GET("/watch/:name")
				services.GET("/status/:name")
				services.PUT("/status/:name")
			}
			deployments := v1.Group("/deployments")
			{
				deployments.GET("/get")
				deployments.GET("/get/:name")
				deployments.POST("/create")
				deployments.PUT("/create/:name")
				deployments.DELETE("/delete")
				deployments.DELETE("/delete/:name")
				deployments.GET("/watch")
				deployments.GET("/watch/:name")
				deployments.GET("/status/:name")
				deployments.PUT("/status/:name")
			}
			hpas := v1.Group("/hpa")
			{
				hpas.GET("/get")
				hpas.GET("/get/:name")
				hpas.POST("/create")
				hpas.PUT("/create/:name")
				hpas.DELETE("/delete")
				hpas.DELETE("/delete/:name")
				hpas.GET("/watch")
				hpas.GET("/watch/:name")
				hpas.GET("/status/:name")
				hpas.PUT("/status/:name")
			}
			nodes := v1.Group("/nodes")
			{
				nodes.GET("/get")
				nodes.GET("/get/:name")
				nodes.POST("/create")
				nodes.PUT("/create/:name")
				nodes.DELETE("/delete")
				nodes.DELETE("/delete/:name")
				nodes.GET("/watch")
				nodes.GET("/watch/:name")
				nodes.GET("/status/:name")
				nodes.PUT("/status/:name")
			}
		}
	}
}
