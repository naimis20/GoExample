package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const perPage = 15
const bookingURL = "https://www.booking.com/searchresults.en-us.html?label=gen173nr-1DCAEoggI46AdIM1gEaFyIAQGYATG4ARfIAQzYAQPoAQH4AQKIAgGoAgO4AtC_o-sFwAIB&sid=8d040eb4c0310711e3b054c9c33fedb6&sb=1&src=index&src_elem=sb&error_url=https%3A%2F%2Fwww.booking.com%2Findex.html%3Flabel%3Dgen173nr-1DCAEoggI46AdIM1gEaFyIAQGYATG4ARfIAQzYAQPoAQH4AQKIAgGoAgO4AtC_o-sFwAIB%3Bsid%3D8d040eb4c0310711e3b054c9c33fedb6%3Bsb_price_type%3Dtotal%26%3B&ss=Ioannina&is_ski_area=0&ssne=Ioannina&ssne_untouched=Ioannina&dest_id=-818126&dest_type=city&checkin_year=&checkin_month=&checkout_year=&checkout_month=&group_adults=2&group_children=0&no_rooms=1&b_h4u_keep_filters=&from_sf=1&rows=%d&offset=%d"

//Hotel type
type Hotel struct {
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	RatingAvg string `json:"rating_avg"`
}

func ScrapeBooking(client *http.Client, pages int) (h []Hotel) {
	c := make(chan []Hotel, pages)

	// for i := 0; i < pages; i++ {
	url := fmt.Sprintf(bookingURL, perPage, pages*perPage) //+ fmt.Sprintf("&offset=%d", i*perPage)
	go Scrape(client, url, c)
	// }
	// for i := 0; i < pages; i++ {
	h = append(h, <-c...)
	// }

	return
}

func Scrape(client *http.Client, url string, c chan []Hotel) {
	response, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading http response body", err)
	}

	hotels := []Hotel{}

	document.Find(".sr_item").Each(func(index int, element *goquery.Selection) {
		hotels = append(hotels, Hotel{
			Title:     strings.TrimSpace(element.Find(".sr-hotel__name").Text()),
			Thumbnail: element.Find(".hotel_image").AttrOr("src", ""),
			RatingAvg: strings.TrimSpace(element.Find(".bui-review-score__badge").Text()),
		})
	})

	c <- hotels
}
