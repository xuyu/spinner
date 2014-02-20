package sensor

import (
	"bytes"
	"io/ioutil"
	"syscall"

	"git.code4.in/spinner"
)

func DiskIOCount() (map[string]map[string]int64, error) {
	content, err := ioutil.ReadFile(PROC_PARTITIONS)
	if err != nil {
		return nil, err
	}
	names := make(map[string]bool)
	for _, line := range bytes.Split(content, NEWLINE)[2:] {
		items := bytes.Fields(line)
		name := string(items[len(items)-1])
		names[name] = true
	}
	content, err = ioutil.ReadFile(PROC_DISKSTATS)
	if err != nil {
		return nil, err
	}
	result := make(map[string]map[string]int64)
	for _, line := range bytes.Split(content, NEWLINE) {
		items := bytes.Fields(line)
		name := string(items[2])
		if names[name] {
			result[name] = map[string]int64{
				"reads":  spinner.MustInt64(string(items[3])),
				"rbytes": spinner.MustInt64(string(items[5])),
				"rtime":  spinner.MustInt64(string(items[6])),
				"writes": spinner.MustInt64(string(items[7])),
				"wbytes": spinner.MustInt64(string(items[9])),
				"wtime":  spinner.MustInt64(string(items[10])),
			}
		}
	}
	return result, nil
}

func DiskPartitions() ([][]string, error) {
	content, err := ioutil.ReadFile(PROC_FILESYSTEMS)
	if err != nil {
		return nil, err
	}
	var devs []string
	for _, line := range bytes.Split(content, NEWLINE) {
		if !bytes.HasPrefix(line, NODEV) {
			devs = append(devs, string(bytes.TrimSpace(line)))
		}
	}
	content, err = ioutil.ReadFile(ETC_MTAB)
	if err != nil {
		return nil, err
	}
	var result [][]string
	for _, line := range bytes.Split(content, NEWLINE) {
		items := bytes.Fields(line)
		result = append(result, []string{
			string(items[0]),
			string(items[1]),
			string(items[2]),
			string(items[3]),
		})
	}
	return result, nil
}

func DiskUsage(mountpoint string) (uint64, uint64, error) {
	var fs syscall.Statfs_t
	if err := syscall.Statfs(mountpoint, &fs); err != nil {
		return 0, 0, err
	}
	return fs.Blocks * uint64(fs.Bsize), fs.Bfree * uint64(fs.Bsize), nil
}
