package sensor

import (
	"bytes"
	"fmt"
	"strconv"

	"git.code4.in/spinner"
)

func BootTime() (int64, error) {
	content, err := spinner.ReadOnce(ProcStat)
	if err != nil {
		return 0, err
	}
	for _, line := range bytes.Split(content, Newline) {
		if bytes.HasPrefix(line, BTime) {
			s := string(bytes.TrimSpace(bytes.TrimPrefix(line, BTime)))
			return strconv.ParseInt(s, 10, 64)
		}
	}
	return 0, fmt.Errorf("btime not found")
}
