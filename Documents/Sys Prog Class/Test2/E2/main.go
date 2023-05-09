package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type RequestBody struct {
    Name string `json:"name"`
}

func jsonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var reqBody RequestBody
        err := json.NewDecoder(r.Body).Decode(&reqBody)
        if err != nil {
            http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
            return
        }

        r = r.WithContext(context.WithValue(r.Context(), "requestBody", reqBody))

        next.ServeHTTP(w, r)
    })
}

func helloHandler(w http.ResponseWriter, r *http.Request) {

    reqBody := r.Context().Value("requestBody").(RequestBody)

    log.Printf("Received request with name=%s\n", reqBody.Name)
    w.Write([]byte("Hello, World!"))
}

func main() {
    mux := http.NewServeMux()

    mux.Handle("/hello", jsonMiddleware(http.HandlerFunc(helloHandler)))

    log.Println("Listening on :4000...")
    log.Fatal(http.ListenAndServe(":4000", mux))
}