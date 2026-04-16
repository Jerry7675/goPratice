package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Jerry7675/goPratice/requestsender"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running\n"))
	})
	mux.HandleFunc("/start", startHandler)
	mux.HandleFunc("/status", statusHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	err := requestsender.StartRequester()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
log.Println("Request sender started")
	w.Write([]byte("Request sender started\n"))
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	count := requestsender.RequestsSent()
	w.Write([]byte("Requests sent: " + fmt.Sprint(count) + "\n"))
}
