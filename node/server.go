package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	addr      string
	trustFile string
	hostname  string
)

func init() {
	flag.StringVar(&addr, "http", ":51002", "http server listen address")
	flag.StringVar(&trustFile, "trust", "trust.ips", "trust ips file")

	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err.Error())
	}
	flag.StringVar(&hostname, "hostname", name, "unique hostname in whole datacenter")
}

func internalServerError(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	if err != nil {
		rw.Write([]byte(err.Error()))
	}
}

type authHandler struct {
	trustFile string
	trust     []net.IP
	handler   http.Handler
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

func (a *authHandler) fillTrust() {
	data, err := ioutil.ReadFile(a.trustFile)
	if err != nil {
		if os.IsNotExist(err) {
			a.trust = nil
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
		a.trust = append(a.trust, ip)
	}
}

func (a *authHandler) saveTrust() error {
	ips := []string{}
	if a.trust != nil {
		for _, ip := range a.trust {
			ips = append(ips, ip.String())
		}
	}
	data := strings.Join(ips, "\n")
	return ioutil.WriteFile(a.trustFile, []byte(data), 0600)
}

func (a *authHandler) Trust(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	ip := net.ParseIP(req.FormValue("ip"))
	if ip != nil {
		if a.trust == nil {
			a.trust = []net.IP{}
		}
		a.trust = append(a.trust, ip)
		a.saveTrust()
	}
	b, _ := json.Marshal(a.trust)
	rw.Write(b)
}

func (a *authHandler) Untrust(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	ip := net.ParseIP(req.FormValue("ip"))
	if ip != nil && a.trust != nil {
		slice := []net.IP{}
		for _, trustIP := range a.trust {
			if !trustIP.Equal(ip) {
				slice = append(slice, trustIP)
			}
		}
		a.trust = slice
		a.saveTrust()
	}
	b, _ := json.Marshal(a.trust)
	rw.Write(b)
}

const (
	centralSpinner = "http://central.spinner:51001"
)

func main() {
	flag.Parse()

	go keepAlive()
	go updateRunner()

	auth := &authHandler{
		trustFile: trustFile,
		trust:     []net.IP{},
		handler:   http.DefaultServeMux,
	}
	auth.fillTrust()

	http.HandleFunc("/spinner/node/trust", auth.Trust)
	http.HandleFunc("/spinner/node/untrust", auth.Untrust)
	http.HandleFunc("/spinner/node/terminal", Terminal)
	http.HandleFunc("/spinner/node/edit", Edit)
	http.HandleFunc("/spinner/node/dashboard", Dashboard)
	http.HandleFunc("/spinner/node/filesystem", FileSystem)

	log.Printf("http server listen at [%s]", addr)
	log.Fatal(http.ListenAndServe(addr, auth))
}
