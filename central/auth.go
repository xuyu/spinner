package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	domain string
	https  bool
)

func init() {
	flag.StringVar(&domain, "domain", "", "cookie domain")
	flag.BoolVar(&https, "https", false, "cookie secure")
}

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
	case req.URL.Path == "/login":
		http.ServeFile(rw, req, filepath.Join(staticPath, "login.html"))
	case req.URL.Path == "/session" && req.Method == "POST":
		a.webuiLogin(rw, req)
	default:
		a.serveWebUI(rw, req)
	}
}

const secret = "8a4f8e9538a6e950205358d8ae5b52abfe66b8a8e27aeef907a5e49247679829"

func signature(name string, timestamp string) string {
	h := sha256.New()
	h.Write([]byte(secret))
	h.Write([]byte(name))
	h.Write([]byte(timestamp))
	return hex.EncodeToString(h.Sum(nil))
}

func verifySession(s string) bool {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return false
	}
	items := strings.Split(string(b), "-")
	if len(items) != 3 {
		return false
	}
	return signature(items[2], items[1]) == items[0]
}

func newSession(name string, timestamp string) string {
	sign := signature(name, timestamp)
	value := sign + "-" + timestamp + "-" + name
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func gotoLogin(rw http.ResponseWriter, req *http.Request) {
	http.Redirect(rw, req, "/login", http.StatusFound)
}

const SessionName = "SPINNER"

func (a *authHandler) serveWebUI(rw http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(SessionName)
	if err != nil {
		gotoLogin(rw, req)
		return
	}
	if !verifySession(cookie.Value) {
		gotoLogin(rw, req)
		return
	}
	a.webuiHandler.ServeHTTP(rw, req)
}

func (a *authHandler) webuiLogin(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	name := req.PostFormValue("name")
	password := req.PostFormValue("password")
	if name == "admin" && password == "123" {
		now := time.Now()
		exp := now.AddDate(0, 0, 7)
		timestamp := strconv.FormatInt(now.Unix(), 10)
		cookie := &http.Cookie{
			Name:     SessionName,
			Value:    newSession(name, timestamp),
			Path:     "/",
			Domain:   domain,
			Expires:  exp,
			MaxAge:   86400 * 7,
			HttpOnly: true,
			Secure:   https,
		}
		http.SetCookie(rw, cookie)
		http.Redirect(rw, req, "/", http.StatusFound)
		return
	}
	gotoLogin(rw, req)
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
