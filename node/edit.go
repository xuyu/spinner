package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"git.code4.in/spinner"
)

func open(rw http.ResponseWriter, req *http.Request, file string) {
	filecontent, err := ioutil.ReadFile(file)
	if err != nil {
		internalServerError(rw, err)
		log.Printf("open file [%s] error: %s", file, err.Error())
		return
	}
	rw.Write(filecontent)
	log.Printf("open file [%s]", file)
}

func readBody(rw http.ResponseWriter, req *http.Request) []byte {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		internalServerError(rw, err)
		log.Printf("read request body error: %s", err.Error())
		return nil
	}
	defer req.Body.Close()
	return body
}

func absLink(link string) (string, error) {
	p, err := os.Readlink(link)
	if err != nil {
		log.Printf("abs link [%s] error: %s", link, err.Error())
		return p, err
	}
	if filepath.IsAbs(p) {
		return p, nil
	}
	return filepath.Join(link, p), nil
}

func save(rw http.ResponseWriter, req *http.Request, file string) {
	body := readBody(rw, req)
	if body == nil {
		return
	}
	mode, err := spinner.FilePerm(file, 0644)
	if err != nil && !os.IsNotExist(err) {
		internalServerError(rw, err)
		log.Printf("stat file [%s] error: %s", file, err.Error())
		return
	}
	switch {
	case mode&os.ModeSymlink != 0:
		file, err = absLink(file)
		if err != nil {
			internalServerError(rw, err)
			return
		}
	case mode.IsRegular():
	default:
		internalServerError(rw, fmt.Errorf("invalid file: %s", file))
		log.Printf("invalid file mode: %s", file)
		return
	}
	if err := ioutil.WriteFile(file, body, mode.Perm()); err != nil {
		internalServerError(rw, err)
		log.Printf("write file [%s] error: %s", file, err.Error())
		return
	}
	log.Printf("save file [%s]", file)
}

// Edit handles http request
// GET method returns the file's content
// POST method save the body as file's content
func Edit(rw http.ResponseWriter, req *http.Request) {
	file := req.URL.Query().Get("file")
	if file == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	switch strings.ToUpper(req.Method) {
	case "GET":
		open(rw, req, file)
	case "POST":
		save(rw, req, file)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("http request method not allowed: %s", req.Method)
	}
}
