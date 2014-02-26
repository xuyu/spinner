package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func Filesystem(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	p := req.FormValue("path")
	info, err := os.Lstat(p)
	if err != nil {
		log.Println(err.Error())
		if os.IsNotExist(err) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := make(map[string]interface{})
	switch {
	case info.IsDir():
		fs, err := ioutil.ReadDir(p)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		result["type"] = "dir"
		if len(fs) > 1024 {
			result["length"] = len(fs)
			break
		}
		files := []string{}
		for _, f := range fs {
			name := f.Name()
			if f.IsDir() {
				name = "+" + name
			} else {
				name = "-" + name
			}
			files = append(files, name)
		}
		result["files"] = files
	case info.Mode()&os.ModeSymlink != 0:
		result["type"] = "link"
		link, _ := os.Readlink(p)
		if !filepath.IsAbs(link) {
			link = filepath.Join(p, link)
		}
		result["link"] = link
	default:
		result["type"] = "other"
		result["mtime"] = info.ModTime()
		result["size"] = info.Size()
	}
	b, err := json.Marshal(result)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	rw.Write(b)
}
