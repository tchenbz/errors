package main

import (
	"log"
	"net/http"
)

//write middleware
func firstMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		//this is executed on the way down to the handler
		log.Println("Executing firstMiddleware")
		next.ServeHTTP(w, r)
		//this is executed on the way up to the client
		log.Println("Executing firstMiddleware again")
	})
}

func secondMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		//this is executed on the way down to the handler
		log.Println("Executing secondMiddleware")
		if r.URL.Path == "/second" {
			return
		}
		next.ServeHTTP(w, r)
		//this is executed on the way up to the client
		log.Println("Executing secondMiddleware again")
	})
}

//create a handler function
func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing a handler...")
	w.Write([]byte("Hello World!"))
}


func main() {
	//multiplexer(map)/router
	mux := http.NewServeMux()
	mux.Handle("/verify", firstMiddleware(secondMiddleware(http.HandlerFunc(helloHandler))))

	log.Print("starting server on :4000")
	err := http.ListenAndServe(":4000", mux) //create a server
	log.Fatal(err)
}