package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	keepAliveDuration int
)

func init() {
	flag.IntVar(&keepAliveDuration, "keepalive", 180, "duration send keepalive with central server")
}

func keepAlive() {
	u := fmt.Sprintf("%s/spinner/central/keepalive?hostname=%s", centralSpinner, url.QueryEscape(hostname))
	for {
		resp, err := http.Get(u)
		if err != nil {
			log.Println(err.Error())
		} else if resp.StatusCode != http.StatusOK {
			log.Printf("%s %s", u, resp.Status)
		}
		time.Sleep(keepAliveDuration * time.Second)
	}
}
