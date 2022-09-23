package service

import (
	"encoding/json"
	"fmt"
	"github.com/BHAVYAghub/Youtube-API/driver"
	log "github.com/BHAVYAghub/Youtube-API/logging"
	"github.com/BHAVYAghub/Youtube-API/models"
	"io"
	"net/http"
	"time"
)

type YT struct {
	baseURL              string
	query                string
	key                  string
	firstRecordStartTime time.Time
	mongo                driver.MongoDriver
}

func New(baseURL string, query string, key string) YT {
	return YT{baseURL: baseURL, query: query, key: key}
}

func (yt YT) GetVideoDetails(time time.Time) []models.Item {
	client := &http.Client{}

	url := yt.baseURL + "/search"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error(err.Error())
	}

	// set query params to be sent.
	q := req.URL.Query()
	q.Add("part", "snippet")
	q.Add("q", yt.query)
	q.Add("key", yt.key)
	q.Add("type", "video")
	q.Add("order", "date")
	q.Add("publishedAfter", time.UTC().String())

	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	items := models.Items{}
	err = json.Unmarshal(body, &items)
	if err != nil {
		log.Info(err.Error())
		return nil
	}
	return items.Item
}

func (yt YT) FetchAndInsertRecords() error {
	t, err := yt.mongo.GetLastRecordTime()
	if err != nil {
		return err
	}

	if t == nil {
		log.Info("Fetching first record.")
		t = &yt.firstRecordStartTime
	}

	items := yt.GetVideoDetails(*t)

	dbRecords := make([]models.YTRecord, 0)
	for i := range items {
		record := models.YTRecord{}
		record.PublishedAt = items[i].Snippet.PublishedAt
		record.VideoTitle = items[i].Snippet.Title
		record.Description = items[i].Snippet.Desc

		dbRecords = append(dbRecords, record)
	}

	err = yt.mongo.SaveAll(dbRecords)
	if err != nil {
		return err
	}

	return nil
}
