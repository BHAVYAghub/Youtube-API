package external

import (
	"github.com/BHAVYAghub/Youtube-API/models/service"
	"time"
)

type YouTube interface {

	// GetVideoDetails calls the youtube service and returns the response from it.
	GetVideoDetails(time time.Time, pageToken string) *service.YTResponse
}
