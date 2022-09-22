package models

type Response struct {
	page         int
	limit        int
	totalRecords int
	YTRecords    []YTRecord
}
