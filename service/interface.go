package service

import (
	srvcModel "github.com/BHAVYAghub/Youtube-API/models/service"
)

type Service interface {
	// GetData returns the paginated youtube video data stored in DB.
	GetData(limit, page int) (*srvcModel.GetResponse, *srvcModel.SvcError)

	// GetSearchResult returns the paginated youtube video data stored in DB on the basis of search query.
	GetSearchResult(limit, page int, search string) (*srvcModel.GetResponse, *srvcModel.SvcError)

	// FetchAndInsertRecords Fetches record from external service & saves it in DB.
	FetchAndInsertRecords() error
}
