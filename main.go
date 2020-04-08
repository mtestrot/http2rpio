package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}
	defer rpio.Close()

	socket, err := net.Listen("tcp", "localhost:8090")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listening on http://%s\n", socket.Addr().String())

	http.HandleFunc("/pins/", func(w http.ResponseWriter, r *http.Request) {
		// http://localhost:8090/pins/17?status=on
		pinNr, parsePinErr := strconv.Atoi(r.URL.Path[len("/pins/"):])
		if parsePinErr != nil {
			log.Printf("an unexpected error occured: %v", parsePinErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		queryParams, parseQueryErr := url.ParseQuery(r.URL.RawQuery)
		if parseQueryErr != nil {
			log.Printf("an unexpected error occured: %v", parseQueryErr)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var status string
		if len(queryParams["status"]) == 0 {
			log.Println("'status' query parameter is missing")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		status = queryParams["status"][0]

		if !strings.EqualFold(status, "on") && !strings.EqualFold(status, "off") {
			log.Printf("'status' query parameter must be 'on' or 'off' but was '%v'\n", status)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		fmt.Printf("turn pin %v %v\n", pinNr, status)
		pin := rpio.Pin(pinNr)
		pin.Output()
		if strings.EqualFold(status, "on") {
			pin.Low()
		} else {
			pin.High()
		}
	})

	log.Fatal(http.Serve(socket, nil))
}
