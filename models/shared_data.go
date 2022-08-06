package models

// SharedData is data found in instagram page
type SharedData struct {
	EntryData struct {
		TagPage []struct {
			GraphQL struct {
				HashTag struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Media struct {
						Count int        `json:"count"`
						Edges []DataEdge `json:"edges"`
					} `json:"edge_hashtag_to_media"`
					TopPosts struct {
						Edges []DataEdge `json:"edges"`
					}
				} `json:"hashtag"`
			} `json:"graphql"`
		} `json:"TagPage"`
	} `json:"entry_data"`
}

func (s *SharedData) ToInstagramPosts(maxResult int64) []InstagramPost {
	posts := []InstagramPost{}
	if len(s.EntryData.TagPage) > 0 {
		var edges []DataEdge
		if maxResult < int64(len(s.EntryData.TagPage[0].GraphQL.HashTag.Media.Edges)) {
			edges = s.EntryData.TagPage[0].GraphQL.HashTag.Media.Edges[0:maxResult]
		} else {
			edges = s.EntryData.TagPage[0].GraphQL.HashTag.Media.Edges
		}
		for _, edge := range edges {
			posts = append(posts, InstagramPost{
				ImageURL:     edge.Node.DisplayURL,
				ThumbnailURL: edge.Node.ThumbnailURL,
				ShortCode:    edge.Node.ShortCode,
				Width:        edge.Node.Dimensions.Width,
				Height:       edge.Node.Dimensions.Height,
			})
		}
	}
	return posts
}
