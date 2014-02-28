package main

import (
	"net"
	"net/http"
	"time"
)

func keepAlive(rw http.ResponseWriter, req *http.Request) {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.ParseForm()
	hostname := req.FormValue("hostname")
	if hostname == "" {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	ms := datacenter.findAllMachines(hostname)
	if ms == nil || len(ms) == 0 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	timestamp := time.Now().Unix()
	for _, m := range ms {
		m.KeepAlive = timestamp
		m.IP = host
	}
}
