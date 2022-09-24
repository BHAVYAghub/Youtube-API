package external

import (
	"github.com/BHAVYAghub/Youtube-API/models/service"
	"time"
)

type YouTube interface {
	GetVideoDetails(time time.Time, pageToken string) *service.YTResponse
}
