package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/naimis20/GoExample/scraper"
)

func Serve(h []scraper.Hotel) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("hotels.html")
		if err != nil {
			log.Fatal(err.Error())
		}
		t.Execute(w, h)
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(h)
	})

	http.ListenAndServe(":8080", nil)

}
