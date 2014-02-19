package sensor

import (
	"bytes"
	"io/ioutil"
	"strconv"
)

func NetIOCount() (map[string]map[string]int64, error) {
	content, err := ioutil.ReadFile("/pro/net/dev")
	if err != nil {
		return nil, err
	}
	result := make(map[string]map[string]int64)
	for _, line := range bytes.Split(content, []byte("\n"))[2:] {
		i := bytes.LastIndex(line, []byte(":"))
		name := string(line[:i])
		fields := bytes.Fields(bytes.TrimSpace(line[i+1:]))
		rbytes, err := strconv.ParseInt(string(fields[0]), 10, 64)
		if err != nil {
			return nil, err
		}
		rpackets, err := strconv.ParseInt(string(fields[1]), 10, 64)
		if err != nil {
			return nil, err
		}
		errin, err := strconv.ParseInt(string(fields[2]), 10, 64)
		if err != nil {
			return nil, err
		}
		dropin, err := strconv.ParseInt(string(fields[3]), 10, 64)
		if err != nil {
			return nil, err
		}
		sbytes, err := strconv.ParseInt(string(fields[8]), 10, 64)
		if err != nil {
			return nil, err
		}
		spackets, err := strconv.ParseInt(string(fields[9]), 10, 64)
		if err != nil {
			return nil, err
		}
		errout, err := strconv.ParseInt(string(fields[10]), 10, 64)
		if err != nil {
			return nil, err
		}
		dropout, err := strconv.ParseInt(string(fields[11]), 10, 64)
		if err != nil {
			return nil, err
		}
		result[name] = map[string]int64{
			"rbytes":   rbytes,
			"rpackets": rpackets,
			"errin":    errin,
			"dropin":   dropin,
			"sbytes":   sbytes,
			"spackets": spackets,
			"errout":   errout,
			"dropout":  dropout,
		}
	}
	return result, nil
}
