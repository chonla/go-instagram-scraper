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

type TagScraper struct {
	userAgent string
}

func NewTag(ua string) *TagScraper {
	return &TagScraper{
		userAgent: ua,
	}
}

func (t *TagScraper) Scrape(tag string, maxResult int64) ([]models.InstagramPost, error) {
	options := []func(*colly.Collector){}

	if t.userAgent != "" {
		options = append(options, colly.UserAgent(t.userAgent))
	}

	c := colly.NewCollector(options...)

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
		if sharedData == nil {
			sharedDataIndex := strings.Index(e.Text, "window._sharedData = ")
			if sharedDataIndex > -1 {
				logrus.Debug("Shared data found")
				sharedDataText := e.Text[sharedDataIndex+21 : len(e.Text)-1]
				err = json.Unmarshal([]byte(sharedDataText), &sharedData)
				if err != nil {
					logrus.Debugf("Unable to unmarshal JSON %s", sharedDataText)
				} else {
					logrus.Debug("This is what I've got")
					b, _ := json.Marshal(sharedData)
					logrus.Debug(string(b))
				}
			} else {
				logrus.Debug("Shared data not found in the following context")
				logrus.Debug(e.Text)
			}
		} else {
			logrus.Debug("Shared data has been found, skipped ...")
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
