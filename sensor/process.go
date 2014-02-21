package sensor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"

	"git.code4.in/spinner"
)

func ProcessPids() ([]int, error) {
	infos, err := ioutil.ReadDir(PROC)
	if err != nil {
		return nil, err
	}
	var pids []int
	for _, info := range infos {
		if info.IsDir() && unicode.IsDigit(rune(info.Name()[0])) {
			pid, err := strconv.Atoi(info.Name())
			if err == nil {
				pids = append(pids, pid)
			}
		}
	}
	return pids, nil
}

type Process struct {
	Pid     int
	Name    string
	Exe     string
	Cmdline []string
}

func (p *Process) GetName() (string, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", p.Pid))
	if err != nil {
		return "", err
	}
	start := bytes.Index(content, []byte("("))
	end := bytes.LastIndex(content, []byte(")"))
	return string(content[start+1 : end]), nil
}

func (p *Process) GetExe() (string, error) {
	link, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", p.Pid))
	if err != nil {
		return link, err
	}
	return strings.Split(strings.Fields(link)[0], "\x00")[0], nil
}

func (p *Process) GetCmdline() ([]string, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", p.Pid))
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimRight(string(content), "\x00"), "\x00"), nil
}

func (p *Process) GetUid() (string, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/status", p.Pid))
	if err != nil {
		return "", err
	}
	for _, line := range bytes.Split(content, NEWLINE) {
		if bytes.HasPrefix(line, []byte("Uid")) {
			i := bytes.Index(line, []byte(":"))
			n := bytes.Fields(bytes.TrimSpace(line[i+1:]))[0]
			return string(n), nil
		}
	}
	return "", fmt.Errorf("get process[%d] uid failed", p.Pid)
}

func (p *Process) GetFds() (int, error) {
	fds, err := ioutil.ReadDir(fmt.Sprintf("/proc/%d/fd", p.Pid))
	if err != nil {
		return 0, err
	}
	return len(fds), nil
}

func (p *Process) GetCwd() (string, error) {
	return os.Readlink(fmt.Sprintf("/proc/%d/cwd", p.Pid))
}

func (p *Process) GetThreadsNum() (int, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/status", p.Pid))
	if err != nil {
		return 0, err
	}
	for _, line := range bytes.Split(content, NEWLINE) {
		if bytes.HasPrefix(line, []byte("Threads")) {
			n := bytes.TrimSpace(bytes.Split(line, []byte(":"))[1])
			return spinner.MustInt(string(n)), nil
		}
	}
	return 0, fmt.Errorf("get process[%d] threads failed", p.Pid)
}

func (p *Process) GetSocketConnectionsNum() (int, error) {
	fds, err := ioutil.ReadDir(fmt.Sprintf("/proc/%d/fd", p.Pid))
	if err != nil {
		return 0, err
	}
	var count int
	for _, fd := range fds {
		name := fmt.Sprintf("/proc/%d/fd/%s", p.Pid, fd.Name())
		link, err := os.Readlink(name)
		if err == nil && strings.HasPrefix(link, "socket:[") {
			count += 1
		}
	}
	return count, nil
}

func (p *Process) GetState() (string, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/status", p.Pid))
	if err != nil {
		return "", err
	}
	for _, line := range bytes.Split(content, NEWLINE) {
		if bytes.HasPrefix(line, []byte("State")) {
			s := bytes.Split(line, []byte(":"))[1]
			return string(bytes.TrimSpace(s)), nil
		}
	}
	return "", fmt.Errorf("get process[%d] state failed", p.Pid)
}

func (p *Process) GetMemInfo() (int64, int64, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/statm", p.Pid))
	if err != nil {
		return 0, 0, err
	}
	items := bytes.Fields(bytes.TrimSpace(content))
	pagesize := int64(os.Getpagesize())
	vm := spinner.MustInt64(string(items[0]))
	mem := spinner.MustInt64(string(items[1]))
	return pagesize * vm, pagesize * mem, nil
}
