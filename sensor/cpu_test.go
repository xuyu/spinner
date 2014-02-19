package sensor

import (
	"fmt"
	"runtime"
	"testing"
)

func TestCPUCountLogical(t *testing.T) {
	n, err := CPUCountLogical()
	if err != nil {
		t.Error(err)
	}
	if runtime.NumCPU() != n {
		t.Fail()
	}
}

func TestCPUCountPhysical(t *testing.T) {
	n, err := CPUCountPhysical()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("System Physical CPU: %d\n", n)
}

func TestCPUTimes(t *testing.T) {
	times, err := CPUTimes()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("System CPU Times: %#v\n", times)
}
