package sensor

import (
	"bytes"
	"io/ioutil"
	"strconv"

	"git.code4.in/spinner"
)

func CPUCountLogical() (int, error) {
	content, err := spinner.ReadOnce("/proc/stat")
	if err != nil {
		return 0, err
	}
	num := bytes.Count(content, []byte("cpu")) - 1
	return num, nil
}

func CPUCountPhysical() (int, error) {
	content, err := spinner.ReadOnce("/proc/cpuinfo")
	if err != nil {
		return 0, err
	}
	found := make(map[string]bool)
	for _, line := range bytes.Split(content, []byte("\n")) {
		if bytes.HasPrefix(line, []byte("physical id")) {
			found[string(line)] = true
		}
	}
	return len(found), nil
}

func CPUTimes() (map[string]float64, error) {
	content, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}
	line := bytes.Split(content, []byte("\n"))[0]
	items := bytes.Fields(line)[1:]
	total := float64(0.0)
	for _, item := range items {
		i, err := strconv.Atoi(string(item))
		if err != nil {
			return nil, err
		}
		total += float64(i)
	}
	user, err := strconv.Atoi(string(items[0]))
	if err != nil {
		return nil, err
	}
	nice, err := strconv.Atoi(string(items[1]))
	if err != nil {
		return nil, err
	}
	system, err := strconv.Atoi(string(items[2]))
	if err != nil {
		return nil, err
	}
	iowait, err := strconv.Atoi(string(items[4]))
	if err != nil {
		return nil, err
	}
	return map[string]float64{
		"user":   float64(user) / total * 100,
		"nice":   float64(nice) / total * 100,
		"system": float64(system) / total * 100,
		"iowait": float64(iowait) / total * 100,
	}, nil
}
