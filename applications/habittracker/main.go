package main

import (
	"log"
	"net/http"

	"github.com/dan-solli/homeapps/habittracker/templates"
)

const addr = ":6001"

func main() {
	templates, err := templates.NewTemplates()
	if err != nil {
		log.Fatal(err)
	}

	handler, err := trackerhttp.NewTrackerHandler(views.NewTrackerView(templates), views.NewIndexView(templates))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening to port %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
