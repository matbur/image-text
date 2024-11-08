package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/matbur/image-text/server"
)

type config struct {
	Addr string `envconfig:"ADDR" default:":8021"`
	Mode string `envconfig:"MODE"`
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

	switch cfg.Mode {
	case "TEST":
		mode2(cfg.Addr)
	default:
		mode1(cfg.Addr)
	}
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

	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}
	u := fmt.Sprintf("http://%s/3000x200/steel_blue/yellow?text=abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", addr)
	resp, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}

	body := resp.Body
	defer body.Close()

	bb, err := io.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("image.png", bb, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
