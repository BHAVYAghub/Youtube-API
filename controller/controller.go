package controller

import (
	// TODO: fmt and imports
	"errors"
	"fmt"
	log "github.com/BHAVYAghub/Youtube-API/logging"
	models "github.com/BHAVYAghub/Youtube-API/models/controller"
	"github.com/BHAVYAghub/Youtube-API/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Controller struct {
	srvc service.Service
}

func New(srvc service.Service) *Controller {
	return &Controller{srvc: srvc}
}

func (ctrl *Controller) GetData(c *gin.Context) {
	l, p, err := getLimitAndPage(c.Request)
	if err != nil {
		p = 1
		l = 5
	}

	resp, srvcErr := ctrl.srvc.GetData(l, p)
	if srvcErr != nil && srvcErr.Message != "" {
		log.Error("Failed to process GetData request", zap.String("ErrorMessage", srvcErr.Message),
			zap.Int("Status", srvcErr.ResponseCode))
		c.IndentedJSON(srvcErr.ResponseCode, &models.ErrorResponse{ErrorMessage: srvcErr.Message})
		return
	}

	c.IndentedJSON(http.StatusOK, &models.Response{
		PageDetails: models.PageDetails{
			Page:         p,
			Limit:        l,
			TotalRecords: resp.Count},
		YTRecords: resp.Record,
	})
}

func (ctrl *Controller) GetSearchData(c *gin.Context) {
	l, p, err := getLimitAndPage(c.Request)
	if err != nil {
		p = 1
		l = 5
	}

	search := c.Request.URL.Query().Get("search")
	if search == "" {
		errMsg := fmt.Sprint("[search] : cannot be null or empty")
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorMessage: errMsg,
		})
		return
	}

	resp, srvcErr := ctrl.srvc.GetSearchResult(l, p, search)
	if srvcErr != nil && srvcErr.Message != "" {
		log.Error("Failed to process GetSearchResult request", zap.String("ErrorMessage", srvcErr.Message),
			zap.Int("Status", srvcErr.ResponseCode))
		c.IndentedJSON(srvcErr.ResponseCode, &models.ErrorResponse{ErrorMessage: srvcErr.Message})
		return
	}

	c.IndentedJSON(http.StatusOK, &models.Response{
		PageDetails: models.PageDetails{
			Page:         p,
			Limit:        l,
			TotalRecords: resp.Count},
		YTRecords: resp.Record,
	})
}

func getLimitAndPage(request *http.Request) (int, int, error) {
	limit := request.URL.Query().Get("limit")
	if limit == "" {
		errMsg := fmt.Sprint("[limit] : cannot be null or empty")
		return 0, 0, errors.New(errMsg)
	}

	page := request.URL.Query().Get("page")
	if page == "" {
		errMsg := fmt.Sprint("[page] : cannot be null or empty")
		return 0, 0, errors.New(errMsg)
	}

	l, _ := strconv.Atoi(limit)
	p, _ := strconv.Atoi(page)

	return l, p, nil

}
