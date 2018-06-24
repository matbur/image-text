package main

import (
	"io/ioutil"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/matbur/image-text/server"
)

type config struct {
	Addr string `envconfig:"ADDR" default:":8021"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	mode1(cfg.Addr)
}

func mode1(addr string) {
	log.Infof("Start listening at %s", addr)

	http.HandleFunc("/favicon.ico", server.HandleFavicon)
	http.HandleFunc("/", server.HandleMain())
	log.Fatal(http.ListenAndServe(addr, nil))
}

// for debug only
func mode2(addr string) {
	go func() {
		mode1(addr)
	}()

	resp, err := http.Get("http://localhost:8021/3000x200/steel_blue/yellow?text=abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
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
