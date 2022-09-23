package models

type Response struct {
	Page         int        `json:"page"`
	Limit        int        `json:"limit"`
	TotalRecords int64      `json:"totalRecords"`
	YTRecords    []YTRecord `json:"records"`
}
