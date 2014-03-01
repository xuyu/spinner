package main

import (
	"bytes"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// Terminal handles http request which executes shell commands like a linux terminal
func Terminal(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	cmd := strings.TrimSpace(req.FormValue("cmd"))
	if cmd == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	command := exec.Command("/bin/sh", "-c", cmd)
	var buf bytes.Buffer
	command.Stdout = &buf
	command.Stderr = &buf
	if err := command.Start(); err != nil {
		internalServerError(rw, err)
		log.Println(err.Error())
		return
	}
	ch := make(chan error)
	go func(ch chan<- error) {
		ch <- command.Wait()
	}(ch)
	select {
	case <-time.After(time.Minute):
		if p := command.Process; p != nil {
			p.Kill()
		}
		rw.WriteHeader(http.StatusRequestTimeout)
		log.Printf("terminal command [%s] timeout", cmd)
		return
	case err := <-ch:
		if err != nil {
			internalServerError(rw, err)
			log.Printf("terminal command [%s] error: %s", cmd, err.Error())
			return
		}
		rw.Write(buf.Bytes())
		log.Printf("terminal command [%s]", cmd)
	}
}
