package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	addr       string
	dtFile     string
	staticPath string

	datacenter = &DataCenter{Groups: []*Group{}}
)

func init() {
	flag.StringVar(&addr, "addr", ":51001", "http server listen address")
	flag.StringVar(&dtFile, "datacenter", "datacenter.json", "datacenter define data file")
	flag.StringVar(&staticPath, "static", "static", "webui server static path")
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

	webuiHandler := http.NewServeMux()
	webuiHandler.HandleFunc("/spinner/webui/trust", webuiTrust)
	webuiHandler.HandleFunc("/spinner/webui/untrust", webuiUntrust)
	webuiHandler.HandleFunc("/spinner/webui/tree", webuiTree)
	webuiHandler.HandleFunc("/spinner/webui/dashboard", webuiDashboard)
	webuiHandler.HandleFunc("/spinner/webui/terminal", webuiTerminal)
	webuiHandler.HandleFunc("/spinner/webui/filesystem", webuiFilesystem)
	webuiHandler.HandleFunc("/spinner/webui/open", webuiOpen)
	webuiHandler.HandleFunc("/spinner/webui/save", webuiSave)

	staticHandler := http.StripPrefix("/spinner/webui/static/", http.FileServer(http.Dir(staticPath)))

	auth := &authHandler{
		centralHandler: centralHandler,
		webuiHandler:   webuiHandler,
		staticHandler:  staticHandler,
	}

	log.Printf("http server listen at [%s]", addr)
	log.Fatal(http.ListenAndServe(addr, auth))
}
