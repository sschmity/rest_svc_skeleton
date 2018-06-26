package event

import (
	"io"
	"fmt"
	"net/http"
	"log"
)

func InitRestAPI() {
	http.HandleFunc("/createEvent", func(w http.ResponseWriter, r *http.Request) {
		create(w, r)
		fmt.Fprintf(w, "Event Created")
	})
	log.Println("InitRestAPI for event registered")
}

func create(out io.Writer, r *http.Request) {
	r.URL.Query()
	fmt.Fprintln(out, "Creating event", r.URL.Query())
}