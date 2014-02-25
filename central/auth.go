package main

import (
	"log"
	"net"
	"net/http"
)

func isPrivateIP(ip net.IP) bool {
	if ip[0] == 10 {
		return true
	}
	if ip[0] == 172 && ip[1] == 16 {
		return true
	}
	if ip[0] == 192 && ip[1] == 168 {
		return true
	}
	return false
}

type authHandler struct {
	centralHandler *http.ServeMux
}

func (a *authHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/keepalive" {
		host, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}
		ip := net.ParseIP(host)
		if ip == nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !ip.IsLoopback() && !isPrivateIP(ip) {
			rw.WriteHeader(http.StatusNotAcceptable)
			log.Printf("not private ip: %s", ip.String())
			return
		}
		a.centralHandler.ServeHTTP(rw, req)
		return
	}
	rw.WriteHeader(http.StatusForbidden)
}
