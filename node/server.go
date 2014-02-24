package main

import (
	"flag"
	"log"
	"net/http"
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

func main() {
	flag.Parse()

	http.HandleFunc("/spinner/node/terminal", Terminal)
	http.HandleFunc("/spinner/node/edit", Edit)
	http.HandleFunc("/spinner/node/dashboard", Dashboard)

	log.Printf("http server listen at [%s]", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
