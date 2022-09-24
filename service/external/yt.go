package external

import (
	"encoding/json"
	"fmt"
	//log "github.com/BHAVYAghub/Youtube-API/logging"
	"github.com/BHAVYAghub/Youtube-API/models/service"
	"go.uber.org/zap"
	"io"
	"log"
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

func (yt YoutubeSvc) GetVideoDetails(after time.Time, pageToken string) *service.YTResponse {
	client := &http.Client{}

	url := yt.baseURL + "/search"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		// TODO: some places it is fatal and some places it is error.
		log.Println(err.Error())
	}

	// TODO: panic does not give reason and log
	publishedAfter := after.UTC().Format(time.RFC3339)
	if err != nil {
		log.Fatal("Invalid time provided for fetching YT records.", zap.Error(err))
	}

	// set query params to be sent.
	q := req.URL.Query()
	q.Add("part", "snippet")
	q.Add("q", yt.query)
	q.Add("key", yt.key)
	q.Add("type", "video")
	//q.Add("order", "date")
	q.Add("publishedAfter", publishedAfter)

	if pageToken != "" {
		q.Add("pageToken", pageToken)
	}

	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed while reading response from youtube", zap.Error(err))
		return nil
	}

	var response service.YTResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Failed while unmarshalling youtube response", zap.Error(err))
		return nil
	}

	return &response
}
