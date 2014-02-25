package main

import (
	"net"
	"net/http"
	"time"
)

func keepAlive(rw http.ResponseWriter, req *http.Request) {
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
