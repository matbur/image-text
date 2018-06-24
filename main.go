package main

import (
	"io/ioutil"
	"net/http"

	"github.com/matbur/image-text/server"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
}

func main() {
	mode1()
}

func mode1() {
	http.HandleFunc("/favicon.ico", server.HandleFavicon)
	http.HandleFunc("/", server.HandleMain())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mode2() {
	go func() {
		mode1()
	}()

	resp, err := http.Get("http://localhost:8080/300x200/steel_blue/yellow?text=abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if err != nil {
		log.Fatal(err)
	}

	body := resp.Body
	defer body.Close()

	bb, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("image.png", bb, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
