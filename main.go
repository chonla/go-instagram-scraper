package main

import (
	"fmt"
	"instagram-scraper/scraper"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

const defaultPort = "1444"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, use local environment variables instead.")
	}
	tagName := os.Getenv("TAG")
	maxResult, _ := strconv.ParseInt("0"+os.Getenv("MAX_SCRAPED_RESULT"), 10, 64)
	if maxResult <= 0 {
		maxResult = 20
	}
	listenPort := os.Getenv("PORT")
	if listenPort == "" {
		listenPort = defaultPort
	}

	log.Println("Scraped tag:", tagName)

	e := echo.New()
	e.HideBanner = true
	e.GET("/", func(c echo.Context) error {
		tagscraper := scraper.NewTag()
		data, err := tagscraper.Scrape(tagName, maxResult)
		if err != nil {
			c.Echo().Logger.Error(err)
			c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, data)
	})

	fmt.Printf("Intagram tag scraper for hashtag \"%s\"\n", tagName)
	fmt.Printf("Service is listening on :%s\n", listenPort)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", listenPort)))
}
