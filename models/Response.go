package models

type Response struct {
	Page         int
	Limit        int
	TotalRecords int64
	YTRecords    []YTRecord
}
