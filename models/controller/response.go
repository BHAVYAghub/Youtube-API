package controller

import (
	models "github.com/BHAVYAghub/Youtube-API/models/datastore"
)

type Response struct {
	PageDetails PageDetails       `json:"pageDetails"`
	YTRecords   []models.YTRecord `json:"records,omitempty"`
}

type PageDetails struct {
	Page         int   `json:"page,omitempty"`
	Limit        int   `json:"limit,omitempty"`
	TotalRecords int64 `json:"totalRecords,omitempty"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage,omitempty"`
}
