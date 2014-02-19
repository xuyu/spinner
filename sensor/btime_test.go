package sensor

import (
	"fmt"
	"testing"
	"time"
)

func TestBootTime(t *testing.T) {
	timestamp, err := BootTime()
	if err != nil {
		t.Error(err)
	}
	tm := time.Unix(timestamp, 0)
	fmt.Printf("System Boot Time: %s\n", tm.Format(time.RFC3339))
}
