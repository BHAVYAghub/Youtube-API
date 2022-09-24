package service

import (
	srvcModel "github.com/BHAVYAghub/Youtube-API/models/service"
)

type Service interface {
	GetData(int, int) (*srvcModel.GetResponse, *srvcModel.ServiceError)

	GetSearchResult(int, int, string) (*srvcModel.GetResponse, *srvcModel.ServiceError)

	// Fetches record from external service & saves it in db
	FetchAndInsertRecords() error
}
