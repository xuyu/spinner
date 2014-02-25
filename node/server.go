package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "http", ":51002", "http server listen address")
}

func internalServerError(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	if err != nil {
		rw.Write([]byte(err.Error()))
	}
}

type authHandler struct {
	trust   []net.IP
	handler http.Handler
}

func (a *authHandler) ensureTrustIP(ip net.IP) bool {
	if a.trust == nil || ip.IsLoopback() {
		return true
	}
	for _, trustIP := range a.trust {
		if trustIP.Equal(ip) {
			return true
		}
	}
	return false
}

func (a *authHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	ip := net.ParseIP(host)
	if ip == nil || !a.ensureTrustIP(ip) {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	a.handler.ServeHTTP(rw, req)
}

func fillTrustIP(auth *authHandler) {
	data, err := ioutil.ReadFile("trust.ips")
	if err != nil {
		if os.IsNotExist(err) {
			auth.trust = nil
			return
		}
		log.Fatal(err.Error())
	}
	for _, field := range bytes.Fields(data) {
		s := string(field)
		if s == "" {
			continue
		}
		ip := net.ParseIP(s)
		if ip == nil {
			log.Fatalf("invalid trust ip: %s", s)
		}
		auth.trust = append(auth.trust, ip)
	}
}

func main() {
	flag.Parse()

	auth := &authHandler{
		trust:   []net.IP{},
		handler: http.DefaultServeMux,
	}
	fillTrustIP(auth)

	http.HandleFunc("/spinner/node/terminal", Terminal)
	http.HandleFunc("/spinner/node/edit", Edit)
	http.HandleFunc("/spinner/node/dashboard", Dashboard)

	log.Printf("http server listen at [%s]", addr)
	log.Fatal(http.ListenAndServe(addr, auth))
}
