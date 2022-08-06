package main

import (
	"fmt"
	"instagram-scraper/scraper"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

const defaultPort = "1444"

var tags []string = []string{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, use local environment variables instead.")
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

	log.Println("Available tags:", tagNames)

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
				c.Echo().Logger.Error(err)
				c.JSON(http.StatusInternalServerError, err)
			}
			return c.JSON(http.StatusOK, data)
		}
		return c.NoContent(http.StatusNotFound)
	})

	fmt.Printf("Intagram tag scraper for hashtags \"%s\"\n", tagNames)
	fmt.Printf("Service is listening on :%s\n", listenPort)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", listenPort)))
}
