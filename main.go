package main

import (
	"fmt"
	"instagram-scraper/scraper"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

const defaultPort = "1444"

var tags []string = []string{}

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Println("Error loading .env file, use local environment variables instead.")
	}
	tagNames := os.Getenv("TAGS")
	maxResult, _ := strconv.ParseInt("0"+os.Getenv("MAX_SCRAPED_RESULT"), 10, 64)
	if maxResult <= 0 {
		maxResult = 20
	}
	listenPort := os.Getenv("PORT")
	if listenPort == "" {
		listenPort = defaultPort
	}
	logrus.Println("Available tags:", tagNames)

	logFormat := &logrus.TextFormatter{
		DisableLevelTruncation: true,
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
	}
	logrus.SetFormatter(logFormat)
	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logLevel)
	}

	tags = strings.Split(tagNames, ",")

	e := echo.New()
	e.HideBanner = true
	e.GET("/:tag", func(c echo.Context) error {
		tagName := c.Param("tag")

		found := lo.Contains(tags, tagName)

		if found {
			tagscraper := scraper.NewTag()
			data, err := tagscraper.Scrape(tagName, maxResult)
			if err != nil {
				logrus.Error(err)
				c.JSON(http.StatusInternalServerError, err)
			}
			return c.JSON(http.StatusOK, data)
		} else {
			logrus.Debug("Specified tag is not in the allowed tags.")
		}
		return c.NoContent(http.StatusNotFound)
	})

	logrus.Printf("Intagram tag scraper for hashtags \"%s\"\n", tagNames)
	logrus.Printf("Service is listening on :%s\n", listenPort)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", listenPort)))
}
