package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"strconv"
)

func webuiTrust(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	s := req.FormValue("ip")
	ip := net.ParseIP(s)
	if ip == nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if !ip.IsLoopback() && !isPrivateIP(ip) {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}
	result := make(map[string]string)
	for _, m := range datacenter.allMachines() {
		_, code, err := m.TrustCentral(ip.String())
		if err != nil {
			result[m.Hostname] = err.Error()
		} else if code != http.StatusOK {
			result[m.Hostname] = strconv.Itoa(code)
		}
	}
	b, _ := json.Marshal(result)
	rw.Write(b)
}

func webuiUntrust(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	s := req.FormValue("ip")
	result := make(map[string]string)
	for _, m := range datacenter.allMachines() {
		_, code, err := m.UntrustCentral(s)
		if err != nil {
			result[m.Hostname] = err.Error()
		} else if code != http.StatusOK {
			result[m.Hostname] = strconv.Itoa(code)
		}
	}
	b, _ := json.Marshal(result)
	rw.Write(b)
}

func webuiTree(rw http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(datacenter)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	rw.Write(b)
}

func webuiDashboard(rw http.ResponseWriter, req *http.Request) {
	machineDoWebUI(rw, req, "Dashboard")
}

func webuiTerminal(rw http.ResponseWriter, req *http.Request) {
	machineDoWebUI(rw, req, "Terminal", "cmd")
}

func webuiFilesystem(rw http.ResponseWriter, req *http.Request) {
	machineDoWebUI(rw, req, "ListFilesystem", "path")
}

func webuiOpen(rw http.ResponseWriter, req *http.Request) {
	machineDoWebUI(rw, req, "OpenFile", "file")
}

func machineDoWebUI(rw http.ResponseWriter, req *http.Request, name string, args ...string) {
	req.ParseForm()
	hostname := req.FormValue("h")
	m := datacenter.findMachine(hostname)
	if m == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	f := reflect.ValueOf(m).MethodByName(name)
	if f.IsNil() {
		rw.WriteHeader(http.StatusNotImplemented)
		return
	}
	in := []reflect.Value{}
	for _, arg := range args {
		in = append(in, reflect.ValueOf(req.FormValue(arg)))
	}
	vals := f.Call(in)
	if !vals[2].IsNil() {
		rw.WriteHeader(http.StatusInternalServerError)
		err := vals[2].Interface().(error)
		log.Println(err.Error())
		return
	}
	rw.WriteHeader(int(vals[1].Int()))
	rw.Write(vals[0].Bytes())
}

func webuiSave(rw http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	hostname := q.Get("h")
	m := datacenter.findMachine(hostname)
	if m == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	filename := q.Get("file")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	data, code, err := m.SaveFile(filename, body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	rw.WriteHeader(code)
	rw.Write(data)
}
