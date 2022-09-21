package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	// Service Layer
}

func New() *Controller {
	return &Controller{}
}

func (ctrl *Controller) GetData(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, &gin.H{"status": "success"})
}
