package main

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/snake0207/scrap/scrapper"
)

const constFileName string = "jobs.csv"

func main() {
	webStart()
}

func webStart() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}

func handleHome(c echo.Context) error {
	return c.File("./home.html")
}

func handleScrape(c echo.Context) error {
	defer os.Remove(constFileName)
	search := strings.ToLower(c.FormValue("search"))
	scrapper.Scrape(search)
	return c.Attachment(constFileName, "download.csv")
}