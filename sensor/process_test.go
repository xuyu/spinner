package sensor

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"testing"
)

func TestProcessPids(t *testing.T) {
	pids, err := ProcessPids()
	if err != nil {
		t.Error(err)
	}
	if len(pids) == 0 {
		t.Fail()
	}
}

func TestProcessName(t *testing.T) {
	pid := os.Getpid()
	p := &Process{Pid: pid}
	name, err := p.GetName()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(name)
}

func TestProcessExe(t *testing.T) {
	pid := os.Getpid()
	p := &Process{Pid: pid}
	exe, err := p.GetExe()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(exe)
}

func TestProcessCmdline(t *testing.T) {
	pid := os.Getpid()
	p := &Process{Pid: pid}
	cmdline, err := p.GetCmdline()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("[%s]", strings.Join(cmdline, " ")))
}

func TestProcessUid(t *testing.T) {
	pid := os.Getpid()
	p := &Process{Pid: pid}
	uid, err := p.GetUid()
	if err != nil {
		t.Error(err)
	}
	u, err := user.LookupId(uid)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(uid, u.Name)
}

func TestProcessFds(t *testing.T) {
	pid := os.Getpid()
	p := &Process{Pid: pid}
	fds, err := p.GetFds()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fds)
}

func TestProcessCwd(t *testing.T) {
	pid := os.Getpid()
	p := &Process{Pid: pid}
	cwd, err := p.GetCwd()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(cwd)
}

func TestProcessThreadsNum(t *testing.T) {
	pid := os.Getpid()
	p := &Process{Pid: pid}
	num, err := p.GetThreadsNum()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(num)
}

func TestProcessSocketConnectionsNum(t *testing.T) {
	pid := os.Getpid()
	p := &Process{Pid: pid}
	num, err := p.GetSocketConnectionsNum()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(num)
}

func TestProcessState(t *testing.T) {
	pid := os.Getpid()
	p := &Process{Pid: pid}
	state, err := p.GetState()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(state)
}

func TestProcessMemInfo(t *testing.T) {
	pid := os.Getpid()
	p := &Process{Pid: pid}
	vm, mem, err := p.GetMemInfo()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(vm, mem)
}
