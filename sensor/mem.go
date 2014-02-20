package sensor

import (
	"bytes"
	"io/ioutil"

	"git.code4.in/spinner"
)

func MemInfo() (map[string]int64, error) {
	content, err := ioutil.ReadFile(PROC_MEMINFO)
	if err != nil {
		return nil, err
	}
	mem := make(map[string]int64)
	for _, item := range bytes.Split(content, NEWLINE) {
		i := bytes.IndexByte(item, ':')
		if i < 0 {
			continue
		}
		name := string(item[:i])
		right := bytes.TrimSpace(item[i+1:])
		kB := false
		if bytes.HasSuffix(right, MEM_KB) {
			kB = true
			right = bytes.TrimSpace(bytes.TrimSuffix(right, MEM_KB))
		}
		value := spinner.MustInt64(string(right))
		if kB {
			value = value * 1024
		}
		mem[name] = value
	}
	return map[string]int64{
		"total":   mem["MemTotal"],
		"free":    mem["MemFree"],
		"buffers": mem["Buffers"],
		"cached":  mem["Cached"],
		"stotal":  mem["SwapTotal"],
		"sfree":   mem["SwapFree"],
	}, nil
}
