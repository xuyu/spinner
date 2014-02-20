package sensor

import (
	"bytes"
	"io/ioutil"

	"git.code4.in/spinner"
)

func CPUCountLogical() (int, error) {
	content, err := spinner.ReadOnce(PROC_STAT)
	if err != nil {
		return 0, err
	}
	num := bytes.Count(content, []byte("cpu")) - 1
	return num, nil
}

func CPUCountPhysical() (int, error) {
	content, err := spinner.ReadOnce(PROC_CPUINFO)
	if err != nil {
		return 0, err
	}
	found := make(map[string]bool)
	for _, line := range bytes.Split(content, NEWLINE) {
		if bytes.HasPrefix(line, PHYSICAL_ID) {
			found[string(line)] = true
		}
	}
	return len(found), nil
}

func CPUTimes() (map[string]float64, error) {
	content, err := ioutil.ReadFile(PROC_STAT)
	if err != nil {
		return nil, err
	}
	line := bytes.Split(content, NEWLINE)[0]
	items := bytes.Fields(line)[1:]
	total := float64(0.0)
	for _, item := range items {
		total += float64(spinner.MustInt64(string(item)))
	}
	return map[string]float64{
		"user":   float64(spinner.MustInt64(string(items[0]))) / total * 100,
		"nice":   float64(spinner.MustInt64(string(items[1]))) / total * 100,
		"system": float64(spinner.MustInt64(string(items[2]))) / total * 100,
		"iowait": float64(spinner.MustInt64(string(items[3]))) / total * 100,
	}, nil
}
