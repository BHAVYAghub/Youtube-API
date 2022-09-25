package external

import (
	"encoding/json"
	log "github.com/BHAVYAghub/Youtube-API/logging"
	"github.com/BHAVYAghub/Youtube-API/models/service"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type YoutubeSvc struct {
	baseURL string
	query   string
	key     string
}

func NewService(baseURL, query, key string) *YoutubeSvc {
	return &YoutubeSvc{
		baseURL: baseURL,
		query:   query,
		key:     key,
	}
}

func (yt YoutubeSvc) GetVideoDetails(after time.Time, pageToken string) (*service.YTResponse, error) {
	client := &http.Client{}

	url := yt.baseURL + "/search"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error("Error from YT API call. ", zap.Error(err))
		return nil, err
	}

	publishedAfter := after.UTC().Format(time.RFC3339)
	if err != nil {
		log.Error("Invalid time provided for fetching YT records.", zap.Error(err))
		return nil, err
	}

	// set query params to be sent.
	q := req.URL.Query()
	q.Add("part", "snippet")
	q.Add("q", yt.query)
	q.Add("key", yt.key)
	q.Add("type", "video")
	q.Add("publishedAfter", publishedAfter)
	if pageToken != "" {
		q.Add("pageToken", pageToken)
	}

	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error from YT API call. ", zap.Error(err))
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed while reading response from youtube", zap.Error(err))
		return nil, err
	}

	var response service.YTResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error("Failed while unmarshalling youtube response", zap.Error(err))
		return nil, err
	}

	return &response, nil
}
