package service

import (
	"github.com/BHAVYAghub/Youtube-API/datastore"
	log "github.com/BHAVYAghub/Youtube-API/logging"
	storeModel "github.com/BHAVYAghub/Youtube-API/models/datastore"
	srvcModel "github.com/BHAVYAghub/Youtube-API/models/service"
	"github.com/BHAVYAghub/Youtube-API/service/external"

	"net/http"
	"time"

	"go.uber.org/zap"
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

func (yt YTService) GetData(limit, page int) (*srvcModel.GetResponse, *srvcModel.SvcError) {
	count, records, err := yt.mongo.GetAll(limit, page)
	if err != nil {
		return nil, &srvcModel.SvcError{
			Message:      err.Error(),
			ResponseCode: http.StatusInternalServerError,
		}
	}

	if len(records) == 0 || count == 0 {
		log.Error("No Record Found.")
		return nil, &srvcModel.SvcError{
			Message:      "No records found",
			ResponseCode: http.StatusNotFound,
		}
	}

	return &srvcModel.GetResponse{
		Count:  count,
		Record: records,
	}, nil
}

func (yt YTService) GetSearchResult(limit, page int, searchString string) (*srvcModel.GetResponse, *srvcModel.SvcError) {
	count, records, err := yt.mongo.GetByTitleOrDesc(limit, page, searchString)
	if err != nil {
		return nil, &srvcModel.SvcError{
			Message:      err.Error(),
			ResponseCode: http.StatusInternalServerError,
		}
	}

	if records == nil || count == 0 {
		log.Error("No Record Found.")
		return nil, &srvcModel.SvcError{
			Message:      "No records found",
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
		log.Info("Fetching first record", zap.String("time", yt.firstRecordStartTime.String()))
		t = &yt.firstRecordStartTime
	} else {
		log.Info("Fetching records After ", zap.String("time", t.String()))
	}

	pageToken := ""

	dbRecords := make([]storeModel.YTRecord, 0)
	for true {
		log.Info("Calling YoutubeSvc API", zap.String("PageToken", pageToken))
		ytResponse, quotaExceeded, err := yt.ExternalSvc.GetVideoDetails(*t, pageToken)
		if err != nil {
			break
		}

		dbData := transformYoutubeResponse(ytResponse)
		for i := range dbData {
			dbRecords = append(dbRecords, dbData[i])
		}

		pageToken = ytResponse.NextPageToken
		if pageToken == "" && !quotaExceeded {
			break
		}
	}

	// saving records once fetched all pages.
	err = yt.mongo.SaveAll(dbRecords)
	if err != nil {
		return err
	}

	return nil
}

func transformYoutubeResponse(response *srvcModel.YTResponse) []storeModel.YTRecord {
	records := make([]storeModel.YTRecord, 0)

	for _, item := range response.Items {
		record := storeModel.YTRecord{}
		record.PublishedAt = item.Snippet.PublishedAt
		record.VideoTitle = item.Snippet.Title
		record.Description = item.Snippet.Desc

		records = append(records, record)
	}
	return records
}
