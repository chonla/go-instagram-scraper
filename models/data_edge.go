package models

type DataEdge struct {
	Node struct {
		DisplayURL   string `json:"display_url"`
		ThumbnailURL string `json:"thumbnail_src"`
		IsVideo      bool   `json:"is_video"`
		Dimensions   struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"dimensions"`
		ShortCode string `json:"shortcode"`
	}
}
