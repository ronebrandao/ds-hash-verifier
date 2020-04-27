package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/async/", FindHashAsync)
	http.HandleFunc("/sync/", FindHash)

	http.ListenAndServe(":8080", nil)
}
