package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	addr   string
	dtFile string

	datacenter = &DataCenter{Groups: []*Group{}}
)

func init() {
	flag.StringVar(&addr, "addr", ":51001", "http server listen address")
	flag.StringVar(&dtFile, "datacenter", "datacenter.json", "datacenter define data file")
}

func main() {
	flag.Parse()

	if err := datacenter.fill(dtFile); err != nil {
		log.Fatal(err.Error())
	}

	centralHandler := http.NewServeMux()
	centralHandler.HandleFunc("/spinner/central/keepalive", keepAlive)
	centralHandler.HandleFunc("/spinner/central/checkupdate", checkUpdate)
	centralHandler.HandleFunc("/spinner/central/update", doUpdate)

	auth := &authHandler{
		centralHandler: centralHandler,
	}

	log.Printf("http server listen at [%s]", addr)
	log.Fatal(http.ListenAndServe(addr, auth))
}
