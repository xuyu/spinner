package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var (
	terminalShell string
)

func init() {
	flag.StringVar(&terminalShell, "shell", "/bin/bash", "terminal shell")
}

// Terminal handles http request which executes shell commands like a linux terminal
func Terminal(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	cmd := strings.TrimSpace(req.FormValue("cmd"))
	if cmd == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	command := exec.Command(terminalShell, "-c", cmd)
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
			rw.Write([]byte("\n"))
			log.Printf("terminal command [%s] error: %s", cmd, err.Error())
		} else {
			log.Printf("terminal command [%s]", cmd)
		}
		rw.Write(buf.Bytes())
	}
}
