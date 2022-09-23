package service

import (
	srvcModel "github.com/BHAVYAghub/Youtube-API/models/service"
)

type Service interface {
	GetData(int, int) (*srvcModel.GetAllResponse, error)

	// Fetches record from external service & saves it in db
	FetchAndInsertRecords() error
}
