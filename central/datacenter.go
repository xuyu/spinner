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

type Group struct {
	Name     string
	Machines []*Machine
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

func (d *DataCenter) allMachines() map[string]*Machine {
	machines := make(map[string]*Machine)
	for _, gp := range d.Groups {
		for _, m := range gp.Machines {
			machines[m.Hostname] = m
		}
	}
	return machines
}

func (d *DataCenter) whichGroups(m *Machine) []*Group {
	gps := []*Group{}
	for _, gp := range d.Groups {
		for _, ma := range gp.Machines {
			if m.Hostname == ma.Hostname {
				gps = append(gps, gp)
			}
		}
	}
	return gps
}
