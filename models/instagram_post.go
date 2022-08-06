package models

type InstagramPost struct {
	ShortCode    string `json:"shortcode"`
	ImageURL     string `json:"imageUrl"`
	ThumbnailURL string `json:"thumbnailUrl"`
	IsVideo      bool   `json:"isVideo"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}
