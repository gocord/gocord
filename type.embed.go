package gocord

import "time"

type EmbedFooter struct {
}

type EmbedImage struct {
}

type EmbedThumbnail struct {
}

type EmbedVideo struct {
}

type EmbedProvider struct {
}

type EmbedAuthor struct {
}

type EmbedFields struct {
}

type Embed struct {
	Title       string         `json:"title"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	URL         string         `json:"url"`
	Timestamp   time.Time      `json:"timestamp"`
	Color       int            `json:"color"`
	Footer      EmbedFooter    `json:"footer"`
	Image       EmbedImage     `json:"image"`
	Thumbnail   EmbedThumbnail `json:"thumbnail"`
	Video       EmbedVideo     `json:"video"`
	Provider    EmbedProvider  `json:"provider"`
	Author      EmbedAuthor    `json:"author"`
	Fields      []EmbedFields  `json:"fields"`
}
