package service

import "time"

type YTResponse struct {
	PageInfo          PageInfo `json:"pageInfo"`
	PreviousPageToken string   `json:"previousPageToken"`
	NextPageToken     string   `json:"nextPageToken"`
	Items             []Item   `json:"items"`
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type Snippet struct {
	PublishedAt  time.Time `json:"publishedAt"`
	ChannelID    string    `json:"channelId,omitempty"`
	Title        string    `json:"title,omitempty"`
	Desc         string    `json:"description,omitempty"`
	ChannelTitle string    `json:"channelTitle,omitempty"`
}

type Item struct {
	Snippet Snippet `json:"snippet,omitempty"`
}
