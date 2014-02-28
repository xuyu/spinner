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

func readUpdateList() (map[string]map[string]string, error) {
	data, err := ioutil.ReadFile(updateListFile)
	if err != nil {
		return nil, err
	}
	var list map[string]map[string]string
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func genUpdateList(local string, list map[string]map[string]string, hostname string) []string {
	files := []string{local}
	m := datacenter.findMachine(hostname)
	if m == nil {
		return files
	}
	gps := datacenter.whichGroups(m)
	for _, gp := range gps {
		val := list[gp.Name]
		if val != nil && len(val) > 0 {
			for f := range val {
				files = append(files, f)
			}
		}
	}
	return files
}

func checkUpdate(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	version := req.FormValue("version")
	local := spinner.ReadVersion(versionFile)
	if local == "" || version == local {
		rw.WriteHeader(http.StatusNotModified)
		return
	}
	list, err := readUpdateList()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	files := genUpdateList(local, list, req.FormValue("h"))
	b, err := json.Marshal(files)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	rw.Write(b)
}

func findUpdateLocal(list map[string]map[string]string, hostname, remote string) string {
	m := datacenter.findMachine(hostname)
	if m == nil {
		return ""
	}
	gps := datacenter.whichGroups(m)
	var local string
	var name string
	for _, gp := range gps {
		val := list[gp.Name]
		if val != nil && len(val) > 0 {
			for k, v := range val {
				if k == remote && (name == "" || len(gp.Name) > len(name)) {
					name = gp.Name
					local = v
				}
			}
		}
	}
	return local
}

func doUpdate(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	file := req.FormValue("file")
	digest := req.FormValue("md5")
	if file == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	list, err := readUpdateList()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	local := findUpdateLocal(list, req.FormValue("h"), file)
	if local == "" {
		rw.WriteHeader(http.StatusNotFound)
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
	data, err := ioutil.ReadFile(local)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	rw.Write(data)
}
