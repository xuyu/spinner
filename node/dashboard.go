package main

import (
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"strings"
	"syscall"

	"git.code4.in/spinner/sensor"
)

var collectors = [][]interface{}{
	[]interface{}{"hostname", os.Hostname},
	[]interface{}{"btime", sensor.BootTime},
	[]interface{}{"logicpu", sensor.CPUCountLogical},
	[]interface{}{"physicpu", sensor.CPUCountPhysical},
	[]interface{}{"cputimes", sensor.CPUTimes},
	[]interface{}{"diskio", sensor.DiskIOCount},
	[]interface{}{"diskpart", sensor.DiskPartitions},
	[]interface{}{"meminfo", sensor.MemInfo},
	[]interface{}{"netio", sensor.NetIOCount},
	[]interface{}{"diskusage", hardDiskUsage},
	[]interface{}{"load", loadAverage},
}

func hardDiskUsage() (map[string][2]uint64, error) {
	parts, err := sensor.DiskPartitions()
	if err != nil {
		return nil, err
	}
	usages := make(map[string][2]uint64)
	for _, items := range parts {
		if strings.HasPrefix(items[0], "/dev/") {
			mount := items[1]
			total, used, err := sensor.DiskUsage(mount)
			if err != nil {
				return nil, err
			}
			usages[mount] = [2]uint64{total, used}
		}
	}
	return usages, nil
}

func loadAverage() ([3]float64, error) {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	var loads [3]float64
	loads[0] = float64(info.Loads[0]) / float64(65536)
	loads[1] = float64(info.Loads[1]) / float64(65536)
	loads[2] = float64(info.Loads[2]) / float64(65536)
	return loads, err
}

func Dashboard(rw http.ResponseWriter, req *http.Request) {
	d := make(map[string]interface{})
	for _, items := range collectors {
		name := items[0].(string)
		caller := reflect.ValueOf(items[1])
		r := caller.Call(nil)
		if !r[1].IsNil() {
			internalServerError(rw, r[1].Interface().(error))
			return
		}
		d[name] = r[0].Interface()
	}
	b, err := json.Marshal(d)
	if err != nil {
		internalServerError(rw, err)
		return
	}
	rw.Write(b)
}
