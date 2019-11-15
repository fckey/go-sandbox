package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func handler(w http.ResponseWriter, r *http.Request) {
		requestBy := r.Header.Get("X-Goog-Authenticated-User-Email")
	fmt.Fprintf(w, "Request by: %s\n", requestBy)
}

func main() {
	r := chi.NewRouter()

	r.Get("/", handler)

	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}
	http.ListenAndServe(fmt.Sprintf(":%v",port), r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(": %s", port), nil))
}
