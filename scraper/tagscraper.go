package scraper

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"instagram-scraper/models"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

type TagScraper struct{}

func NewTag() *TagScraper {
	return &TagScraper{}
}

func (t *TagScraper) Scrape(tag string, maxResult int64) ([]models.InstagramPost, error) {
	c := colly.NewCollector(
		//colly.CacheDir("./_instagram_cache/"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"),
	)

	var sharedData *models.SharedData
	var err error

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("Referrer", "https://www.instagram.com/explore/tags/"+tag)
		if r.Ctx.Get("gis") != "" {
			gis := fmt.Sprintf("%s:%s", r.Ctx.Get("gis"), r.Ctx.Get("variables"))
			h := md5.New()
			h.Write([]byte(gis))
			gisHash := fmt.Sprintf("%x", h.Sum(nil))
			r.Headers.Set("X-Instagram-GIS", gisHash)
		}
	})

	c.OnHTML("script:not([src])", func(e *colly.HTMLElement) {
		sharedDataIndex := strings.Index(e.Text, "window._sharedData = ")
		if sharedDataIndex > -1 {
			sharedDataText := e.Text[sharedDataIndex+21 : len(e.Text)-1]
			err = json.Unmarshal([]byte(sharedDataText), &sharedData)
		} else {
			logrus.Debug("Shared data not found in the following context")
			logrus.Debug(e.Text)
		}
	})

	var instagramPosts []models.InstagramPost
	err = nil

	c.Visit(fmt.Sprintf("https://www.instagram.com/explore/tags/%s/", tag))

	if sharedData != nil {
		instagramPosts = sharedData.ToInstagramPosts(maxResult)
	}

	return instagramPosts, err
}
