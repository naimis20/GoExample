package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const perPage = 15
const bookingUrl = "https://www.booking.com/searchresults.en-us.html?label=gen173nr-1FCAEoggI46AdIM1gEaFyIAQGYATG4ARfIAQzYAQHoAQH4AQKIAgGoAgO4AvfXmesFwAIB&sid=8d040eb4c0310711e3b054c9c33fedb6&sb=1&src=index&src_elem=sb&error_url=https%3A%2F%2Fwww.booking.com%2Findex.html%3Flabel%3Dgen173nr-1FCAEoggI46AdIM1gEaFyIAQGYATG4ARfIAQzYAQHoAQH4AQKIAgGoAgO4AvfXmesFwAIB%3Bsid%3D8d040eb4c0310711e3b054c9c33fedb6%3Bsb_price_type%3Dtotal%26%3B&ss=Ioannina%2C+Epirus%2C+Greece&is_ski_area=&checkin_year=2019&checkin_month=11&checkin_monthday=6&checkout_year=2019&checkout_month=11&checkout_monthday=9&group_adults=2&group_children=0&no_rooms=1&b_h4u_keep_filters=&from_sf=1&ss_raw=%CE%B9%CF%89%CE%AC%CE%BD%CE%BD%CE%B9%CE%B1&ac_position=1&ac_langcode=en&ac_click_type=b&suggested_term=%CE%B9%CF%89%CE%B1%CE%BD%CE%BD%CE%B9%CE%BD%CE%B1&suggestion_clicked=1&dest_id=-818126&dest_type=city&iata=IOA&place_id_lat=39.663098&place_id_lon=20.852024&search_pageview_id=706e53fb31ce00d3&search_selected=true&search_pageview_id=706e53fb31ce00d3&ac_suggestion_list_length=4&ac_suggestion_theme_list_length=0&order="

type Hotel struct {
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	RatingAvg string `json:"rating_avg"`
}

func ScrapeBooking(client *http.Client, pages int) (h []Hotel) {
	c := make(chan []Hotel, pages)

	for i := 0; i < pages; i++ {
		go Scrape(client, bookingUrl+fmt.Sprintf("&offset=%d", i*perPage), c)
	}
	for i := 0; i < pages; i++ {
		h = append(h, <-c...)
	}

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
