package main

import (
	"log"
	"net/http"
	"os/exec"
	"strings"
)

// Terminal handles http request which executes shell commands like a linux terminal
func Terminal(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	cmd := strings.TrimSpace(req.FormValue("cmd"))
	if cmd == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	items := strings.Fields(cmd)
	output, err := exec.Command(items[0], items[1:]...).CombinedOutput()
	if err != nil {
		internalServerError(rw, err)
		log.Printf("terminal execute command [%s] error: %s", cmd, err.Error())
	} else {
		log.Printf("terminal execute command [%s]", cmd)
	}
	rw.Write(output)
}
