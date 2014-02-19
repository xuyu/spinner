package sensor

import (
	"bytes"
	"io/ioutil"
	"strconv"
)

func MemInfo() (map[string]int64, error) {
	content, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	mem := make(map[string]int64)
	for _, item := range bytes.Split(content, []byte("\n")) {
		i := bytes.IndexByte(item, ':')
		name := string(item[:i])
		right := bytes.TrimSpace(item[i+1:])
		kB := false
		if bytes.HasSuffix(right, []byte("kB")) {
			kB = true
			right = bytes.TrimSpace(bytes.TrimSuffix(right, []byte("kB")))
		}
		value, err := strconv.ParseInt(string(right), 10, 64)
		if err != nil {
			return nil, err
		}
		if kB {
			value = value * 1024
		}
		mem[name] = value
	}
	return mem, nil
}
