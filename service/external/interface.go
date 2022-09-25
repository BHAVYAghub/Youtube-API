package external

import (
	"time"

	"github.com/BHAVYAghub/Youtube-API/models/service"
)

type YouTube interface {

	// GetVideoDetails calls the youtube service and returns the response from it.
	GetVideoDetails(time time.Time, pageToken string) (*service.YTResponse, bool, error)
}
