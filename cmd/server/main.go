package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	http.Handle("/caption", loggingMiddleware(http.HandlerFunc(captionHandler)))
	log.Info("starting server on port 8080")
	http.ListenAndServe(":8080", nil)
}

func captionHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		fmt.Fprintf(w, "GET CAPTION REQUEST")
	case "POST":
		fmt.Fprintf(w, "POST CAPTION REQUEST")
	case "PUT":
		fmt.Fprintf(w, "PUT CAPTION REQUEST")
	default:
		http.Error(w, "Sorry, only GET,PUT and POST methods are supported.", http.StatusMethodNotAllowed)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Infof("uri: %s, method: %s", req.RequestURI, req.Method)
		next.ServeHTTP(w, req)
	})
}
