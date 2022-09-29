package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	latLon := r.Header.Get("X-AppEngine-CityLatLong")
	parts := strings.Split(latLon, ",")
	lat := strings.TrimSpace(parts[0])
	lon := strings.TrimSpace(parts[1])

	ips := strings.Split(r.Header.Get("X-Forwarded-For"), ",")

	userIP := ""
	if len(ips) > 0 {
		userIP = ips[0]
	}

	m := map[string]string{
		"latLong": latLon,
		"city":    r.Header.Get("X-AppEngine-City"),
		"region":  r.Header.Get("X-AppEngine-Region"),
		"country": r.Header.Get("X-AppEngine-Country"),
		"lat":     lat,
		"lon":     lon,
		"ip":      userIP,
	}

	js, err := json.Marshal(m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(js))

}

func main() {
	http.HandleFunc("/", handler)

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	// [END setting_port]
}
