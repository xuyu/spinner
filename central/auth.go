package main

import (
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
)

func isPrivateIP(ip net.IP) bool {
	ip = ip.To4()
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
	centralHandler http.Handler
	webuiHandler   http.Handler
	staticHandler  http.Handler
}

func (a *authHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch {
	case strings.HasPrefix(req.URL.Path, "/spinner/central/"):
		a.serveCentral(rw, req)
	case strings.HasPrefix(req.URL.Path, "/spinner/webui/static/"):
		a.staticHandler.ServeHTTP(rw, req)
	case strings.HasPrefix(req.URL.Path, "/spinner/webui/"):
		a.serveWebUI(rw, req)
	case req.URL.Path == "/":
		http.ServeFile(rw, req, filepath.Join(staticPath, "index.html"))
	default:
		rw.WriteHeader(http.StatusNotFound)
	}
}

func (a *authHandler) serveWebUI(rw http.ResponseWriter, req *http.Request) {
	a.webuiHandler.ServeHTTP(rw, req)
}

func (a *authHandler) serveCentral(rw http.ResponseWriter, req *http.Request) {
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
}
