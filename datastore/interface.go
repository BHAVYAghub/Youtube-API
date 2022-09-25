package datastore

import (
	"time"

	"github.com/BHAVYAghub/Youtube-API/models/datastore"
)

type Database interface {
	// GetAll returns the youtube data in paginated format.
	GetAll(int, int) (int64, []datastore.YTRecord, error)

	// GetByTitleOrDesc returns youtube data based on search query.
	GetByTitleOrDesc(int, int, string) (int64, []datastore.YTRecord, error)

	// GetLastRecordTime returns the time of most recent record present in DB.
	GetLastRecordTime() (*time.Time, error)

	// SaveAll inserts all the youtube records provided.
	SaveAll([]datastore.YTRecord) error
}
