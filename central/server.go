package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
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

func KeepAlive(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	hostname := req.FormValue("hostname")
	if hostname == "" {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	m := datacenter.findMachine(hostname)
	if m == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	m.keepalive = time.Now().Unix()
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	m.host = host
}

func main() {
	flag.Parse()

	data, err := ioutil.ReadFile(dtFile)
	if err != nil {
		log.Printf("read datacenter file error: %s", err.Error())
	} else {
		if err := json.Unmarshal(data, datacenter); err != nil {
			log.Printf("unmarshal datacenter json error: %s", err.Error())
		}
	}

	centralHandler := http.NewServeMux()
	centralHandler.HandleFunc("/keepalive", KeepAlive)

	auth := &authHandler{
		centralHandler: centralHandler,
	}

	log.Printf("http server listen at [%s]", addr)
	log.Fatal(http.ListenAndServe(addr, auth))
}
