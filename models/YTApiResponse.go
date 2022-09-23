package models

import "time"

type Snippet struct {
	PublishedAt time.Time `json:"publishedAt"`
	ChannelID   string    `json:"channelId,omitempty"`
	Title       string    `json:"title,omitempty"`
	Desc        string    `json:"description,omitempty"`
}

type Item struct {
	Snippet Snippet `json:"snippet"`
}

type Items struct {
	Item []Item `json:"items,omitempty"`
}
