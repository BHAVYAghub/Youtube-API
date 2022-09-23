package models

import "time"

type YTRecord struct {
	VideoTitle  string    `json:"video_title,omitempty" bson:"title"`
	Description string    `json:"description,omitempty" bson:"desc"`
	PublishedAt time.Time `json:"publishedAt" bson:"publishedAt"`
}
