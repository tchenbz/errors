package main 

import (
	"log"
	"net/http"

	"github.com/goji/httpauth"
)

func main() {
	authenticationHandler := httpauth.SimpleBasicAuth("tchenbz", "p@ssword")

	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(final)
	mux.Handle("/", authenticationHandler(finalHandler))

	log.Print("Listening on :4000...")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome!"))
}