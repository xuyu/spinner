package sensor

import (
	"bytes"
	"io/ioutil"

	"git.code4.in/spinner"
)

func NetIOCount() (map[string]map[string]int64, error) {
	content, err := ioutil.ReadFile(PROC_NET_DEV)
	if err != nil {
		return nil, err
	}
	result := make(map[string]map[string]int64)
	for _, line := range bytes.Split(content, NEWLINE)[2:] {
		i := bytes.LastIndex(line, []byte(":"))
		name := string(line[:i])
		fields := bytes.Fields(bytes.TrimSpace(line[i+1:]))
		result[name] = map[string]int64{
			"rbytes":   spinner.MustInt64(string(fields[0])),
			"rpackets": spinner.MustInt64(string(fields[1])),
			"errin":    spinner.MustInt64(string(fields[2])),
			"dropin":   spinner.MustInt64(string(fields[3])),
			"sbytes":   spinner.MustInt64(string(fields[8])),
			"spackets": spinner.MustInt64(string(fields[9])),
			"errout":   spinner.MustInt64(string(fields[10])),
			"dropout":  spinner.MustInt64(string(fields[11])),
		}
	}
	return result, nil
}
