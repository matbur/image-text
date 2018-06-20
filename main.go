package main

import (
	"net/http"

	"github.com/matbur/image-text/server"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
}

func main() {
	http.HandleFunc("/", server.Foo())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
