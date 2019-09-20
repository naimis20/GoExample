package main

import (
	"net/http"

	"github.com/naimis20/GoExample/scraper"
	"github.com/naimis20/GoExample/server"
)

func main() {
	pages := 3

	client := &http.Client{}

	server.Serve(scraper.ScrapeBooking(client, pages))
}
