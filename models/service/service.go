package service

import (
	models "github.com/BHAVYAghub/Youtube-API/models/datastore"
)

type GetResponse struct {
	Count  int64
	Record []models.YTRecord
}

type ServiceError struct {
	Message      string
	ResponseCode int
}

// TODO: optimise DB interaction
//		Page Info struct
