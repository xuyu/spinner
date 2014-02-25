package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"git.code4.in/spinner"
)

var (
	versionFile    string
	updateListFile string
)

func init() {
	flag.StringVar(&versionFile, "version", "spinner-central.version", "spinner central version file")
	flag.StringVar(&updateListFile, "update", "node-update.json", "spinner node update list file")
}

func checkUpdate(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	version := req.FormValue("version")
	local := spinner.ReadVersion(versionFile)
	if local == "" || version == local {
		rw.WriteHeader(http.StatusNotModified)
		return
	}
	data, err := ioutil.ReadFile(updateListFile)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	var list map[string]string
	if err := json.Unmarshal(data, &list); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	files := []string{local}
	for name := range list {
		files = append(files, name)
	}
	b, err := json.Marshal(files)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	rw.Write(b)
}

func doUpdate(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	file := req.FormValue("file")
	digest := req.FormValue("md5")
	if file == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadFile(updateListFile)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	var list map[string]string
	if err := json.Unmarshal(data, &list); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	local := list[file]
	if local == "" {
		rw.WriteHeader(http.StatusNotModified)
		return
	}
	s, err := spinner.FileMd5(local)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	if digest == s {
		rw.WriteHeader(http.StatusNotModified)
		return
	}
	data, err = ioutil.ReadFile(local)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	rw.Write(data)
}
