package service

import (
	"github.com/BHAVYAghub/Youtube-API/datastore"
	log "github.com/BHAVYAghub/Youtube-API/logging"
	datastore2 "github.com/BHAVYAghub/Youtube-API/models/datastore"
	srvcModel "github.com/BHAVYAghub/Youtube-API/models/service"
	"github.com/BHAVYAghub/Youtube-API/service/external"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type YTService struct {
	firstRecordStartTime time.Time
	mongo                datastore.Database
	ExternalSvc          external.YouTube
}

func NewYTService(externalSvc external.YouTube, firstRecordStartTime time.Time, mongo datastore.Database) YTService {
	return YTService{
		ExternalSvc:          externalSvc,
		firstRecordStartTime: firstRecordStartTime,
		mongo:                mongo,
	}
}

func (yt YTService) GetData(limit, page int) (*srvcModel.GetResponse, *srvcModel.ServiceError) {
	count, records, err := yt.mongo.GetAll(limit, page)
	if err != nil {
		log.Error("Failed while fetching records for GetAll", zap.Error(err))
		return nil, &srvcModel.ServiceError{
			Message:      err.Error(),
			ResponseCode: http.StatusInternalServerError,
		}
	}

	if records == nil || count == 0 {
		return nil, &srvcModel.ServiceError{
			Message:      "Empty records",
			ResponseCode: http.StatusNotFound,
		}
	}

	return &srvcModel.GetResponse{
		Count:  count,
		Record: records,
	}, nil
}

func (yt YTService) GetSearchResult(limit, page int, searchString string) (*srvcModel.GetResponse, *srvcModel.ServiceError) {
	count, records, err := yt.mongo.GetByTitleOrDesc(limit, page, searchString)
	if err != nil {
		log.Error("Failed while fetching records for GetByTitleOrDesc", zap.Error(err))
		return nil, &srvcModel.ServiceError{
			Message:      err.Error(),
			ResponseCode: http.StatusInternalServerError,
		}
	}

	if records == nil || count == 0 {
		return nil, &srvcModel.ServiceError{
			Message:      "Empty records",
			ResponseCode: http.StatusNotFound,
		}
	}

	return &srvcModel.GetResponse{
		Count:  count,
		Record: records,
	}, nil
}

func (yt YTService) FetchAndInsertRecords() error {
	t, err := yt.mongo.GetLastRecordTime()
	if err != nil {
		return err
	}

	if t == nil {
		log.Info("Fetching first record")
		t = &yt.firstRecordStartTime
	}

	pageToken := ""

	for true {
		// TODO: refactor logs +
		log.Info("Calling YoutubeSvc API for pageToken: " + pageToken)
		ytResponse := yt.ExternalSvc.GetVideoDetails(*t, pageToken)

		log.Info("YoutubeSvc API successfully returned", zap.Any("ResponseBody", ytResponse))

		dbRecords := transformYoutubeResponse(ytResponse)
		log.Info("Saving record in db", zap.Any("Record", dbRecords))

		err = yt.mongo.SaveAll(dbRecords)
		if err != nil {
			return err
		}

		pageToken = ytResponse.NextPageToken

		if pageToken == "" {
			break
		}
	}

	return nil
}

func transformYoutubeResponse(response *srvcModel.YTResponse) []datastore2.YTRecord {
	records := make([]datastore2.YTRecord, 0)

	for _, item := range response.Items {
		record := datastore2.YTRecord{}
		record.PublishedAt = item.Snippet.PublishedAt
		record.VideoTitle = item.Snippet.Title
		record.Description = item.Snippet.Desc

		records = append(records, record)
	}
	return records
}
