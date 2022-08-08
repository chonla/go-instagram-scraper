package scraper

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"instagram-scraper/models"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
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

	if os.Getenv("LOG_LEVEL") == "debug" {
		options = append(options, colly.Debugger(&debug.LogDebugger{}))
	}

	c := colly.NewCollector(options...)

	if os.Getenv("PROXIES") != "" {
		proxies, err := t.randomProxy(strings.Split(os.Getenv("PROXIES"), ",")...)
		if err == nil {
			c.SetProxyFunc(proxies)
			logrus.Debug("Proxies settings are detected and applied")
			logrus.Debug(os.Getenv("PROXIES"))
		} else {
			logrus.Error(err)
		}
	}

	var sharedData *models.SharedData
	var err error

	c.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		// r.Headers.Set("Referrer", "https://www.instagram.com/explore/tags/"+tag)

		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
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
					logrus.Debug("from")
					logrus.Debug(sharedDataText)
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

	scrapedUrl := fmt.Sprintf("https://www.instagram.com/explore/tags/%s/", tag)
	logrus.Debugf("Scraping %s", scrapedUrl)
	c.Visit(scrapedUrl)

	if sharedData != nil {
		instagramPosts = sharedData.ToInstagramPosts(maxResult)
	}

	return instagramPosts, err
}

func (t *TagScraper) randomProxy(urls ...string) (colly.ProxyFunc, error) {
	proxies := []*url.URL{}
	for _, u := range urls {
		parsedU, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, parsedU)
	}

	if len(proxies) > 0 {
		seed := rand.NewSource(time.Now().UnixNano())
		randomizer := rand.New(seed)
		proxI := randomizer.Intn(len(proxies))
		return func(pr *http.Request) (*url.URL, error) {
			logrus.Debug("Scraping target through proxy -> ", proxies[proxI])
			return proxies[proxI], nil
		}, nil
	} else {
		return nil, errors.New("No proxy available")
	}
}
