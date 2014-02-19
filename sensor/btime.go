package sensor

import (
	"bytes"
	"fmt"
	"strconv"

	"git.code4.in/spinner"
)

func BootTime() (int64, error) {
	content, err := spinner.ReadOnce("/proc/stat")
	if err != nil {
		return 0, err
	}
	for _, line := range bytes.Split(content, []byte("\n")) {
		if bytes.HasPrefix(line, []byte("btime")) {
			s := string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("btime"))))
			return strconv.ParseInt(s, 10, 64)
		}
	}
	return 0, fmt.Errorf("btime not found")
}
