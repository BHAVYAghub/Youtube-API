package external

import (
	"fmt"

	log "github.com/BHAVYAghub/Youtube-API/logging"
	"github.com/BHAVYAghub/Youtube-API/models/service"

	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type YoutubeSvc struct {
	baseURL  string
	query    string
	key      []string
	keyIndex int
}

func NewService(baseURL, query string, key []string) *YoutubeSvc {
	return &YoutubeSvc{
		baseURL: baseURL,
		query:   query,
		key:     key,
	}
}

func (yt *YoutubeSvc) GetVideoDetails(after time.Time, pageToken string) (*service.YTResponse, bool, error) {
	client := &http.Client{}

	url := yt.baseURL + "/search"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error("Error from YT API call. ", zap.Error(err))
		return nil, false, err
	}

	publishedAfter := after.UTC().Format(time.RFC3339)
	if err != nil {
		log.Error("Invalid time provided for fetching YT records.", zap.Error(err))
		return nil, false, err
	}

	// set query params to be sent.
	q := req.URL.Query()
	q.Add("part", "snippet")
	q.Add("q", yt.query)
	q.Add("key", yt.key[yt.keyIndex])
	q.Add("type", "video")
	q.Add("publishedAfter", publishedAfter)
	if pageToken != "" {
		q.Add("pageToken", pageToken)
	}

	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error from YT API call. ", zap.Error(err))
		return nil, false, err
	}

	if resp.StatusCode == http.StatusForbidden {
		log.Error("Error from YT API call. Quota exceeded for Key .", zap.String("number", fmt.Sprintf("%d", yt.keyIndex)))
		yt.keyIndex = (yt.keyIndex + 1) % len(yt.key)
		if yt.keyIndex == 0 {
			return &service.YTResponse{NextPageToken: ""}, false, nil
		}
		return &service.YTResponse{NextPageToken: pageToken}, true, nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed while reading response from youtube", zap.Error(err))
		return nil, false, err
	}

	var response service.YTResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error("Failed while unmarshalling youtube response", zap.Error(err))
		return nil, false, err
	}

	return &response, false, nil
}
