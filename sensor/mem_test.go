package sensor

import (
	"fmt"
	"testing"
)

func TestMemInfo(t *testing.T) {
	info, err := MemInfo()
	if err != nil {
		t.Error(err)
	}
	for name, value := range info {
		fmt.Printf("%s: %d\n", name, value)
	}
}
