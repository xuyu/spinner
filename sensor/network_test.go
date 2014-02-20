package sensor

import (
	"fmt"
	"testing"
)

func TestNetIOCount(t *testing.T) {
	result, err := NetIOCount()
	if err != nil {
		t.Error(err)
	}
	for name, data := range result {
		fmt.Printf("%s\n%#v\n", name, data)
	}
}
