package routes

import (
	"github.com/BHAVYAghub/Youtube-API/controller"
	"github.com/gin-gonic/gin"
)

const (
	GET    = "get"
	SEARCH = "search"
)

func AttachRoutes(r *gin.Engine) {
	ctrl := controller.New()
	r.Group("/")
	{
		r.GET(GET, ctrl.GetData)
		r.GET(SEARCH, ctrl.GetData)
	}
}
