package sensor

import (
	"fmt"
	"testing"
)

func TestDiskIOCount(t *testing.T) {
	result, err := DiskIOCount()
	if err != nil {
		t.Error(err)
	}
	for name, data := range result {
		fmt.Printf("%s\n%#v\n", name, data)
	}
}

func TestDiskPartitions(t *testing.T) {
	result, err := DiskPartitions()
	if err != nil {
		t.Error(err)
	}
	for _, items := range result {
		for _, item := range items {
			fmt.Printf("%s ", item)
		}
		fmt.Println()
	}
}

func TestDiskUsage(t *testing.T) {
	all, free, err := DiskUsage("/")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(float64(all-free) / float64(all) * 100)
}
