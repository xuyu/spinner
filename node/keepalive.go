package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func keepAlive() {
	u := fmt.Sprintf("%s/spinner/central/keepalive?hostname=%s", centralSpinner, url.QueryEscape(hostname))
	for {
		resp, err := http.Get(u)
		if err != nil {
			log.Println(err.Error())
		} else if resp.StatusCode != http.StatusOK {
			log.Printf("%s %s", u, resp.Status)
		}
		time.Sleep(3 * time.Minute)
	}
}
