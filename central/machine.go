package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Machine struct {
	Hostname string

	keepalive int64
	host      string
}

func (m *Machine) httpGet(p string, q url.Values) ([]byte, int, error) {
	u := fmt.Sprintf("http://%s:51002%s?%s", m.host, p, q.Encode())
	resp, err := http.Get(u)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

func (m *Machine) Dashboard() ([]byte, int, error) {
	return m.httpGet("/spinner/node/dashboard", url.Values{})
}

func (m *Machine) Terminal(cmd string) ([]byte, int, error) {
	var query = url.Values{}
	query.Add("cmd", cmd)
	return m.httpGet("/spinner/node/terminal", query)
}

func (m *Machine) OpenFile(filename string) ([]byte, int, error) {
	var query = url.Values{}
	query.Add("file", filename)
	return m.httpGet("/spinner/node/edit", query)
}

func (m *Machine) TrustCentral(ip string) ([]byte, int, error) {
	var query = url.Values{}
	query.Add("ip", ip)
	return m.httpGet("/spinner/node/trust", query)
}

func (m *Machine) UntrustCentral(ip string) ([]byte, int, error) {
	var query = url.Values{}
	query.Add("ip", ip)
	return m.httpGet("/spinner/node/untrust", query)
}

func (m *Machine) SaveFile(filename string, content []byte) ([]byte, int, error) {
	u := fmt.Sprintf("http://%s:51002/spinner/node/edit?file=%s", m.host, filename)
	resp, err := http.Post(u, "", bytes.NewReader(content))
	if err != nil {
		return nil, 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}
