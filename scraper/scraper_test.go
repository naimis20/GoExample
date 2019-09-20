package scraper

import (
	"net/http"
	"testing"
)

//Test service is working
func TestScrapeBooking(t *testing.T) {
	client := &http.Client{}

	hotelsResult := ScrapeBooking(client, 1)
	if len(hotelsResult) <= 0 {
		t.Errorf("Service not working")
	}
}
