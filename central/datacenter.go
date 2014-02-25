package main

import (
	"encoding/json"
	"io/ioutil"
)

type DataCenter struct {
	Name     string
	Location string
	Groups   []*Group
}

func (d *DataCenter) findMachine(hostname string) *Machine {
	for _, gp := range d.Groups {
		for _, m := range gp.Machines {
			if m.Hostname == hostname {
				return m
			}
		}
	}
	return nil
}

func (d *DataCenter) fill(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, d)
}
