package driver

import (
	"github.com/BHAVYAghub/Youtube-API/models"
	"time"
)

type Database interface {
	Get(int, int) (int64, []models.YTRecord, error)

	GetByTitleOrDesc(int, int, string) (int64, []models.YTRecord, error)

	GetLastRecordTime() (*time.Time, error)

	SaveAll([]models.YTRecord) error
}
