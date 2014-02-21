package command

import (
	"bytes"
	"os/exec"
)

func SyncCommand(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

func BackgroundCommand(name string, args ...string) error {
	return exec.Command(name, args...).Start()
}

type AsyncCommandResult struct {
	Output []byte
	Error  error
}

func AsyncCommand(name string, args ...string) <-chan AsyncCommandResult {
	ch := make(chan AsyncCommandResult)
	go func(ch chan AsyncCommandResult) {
		cmd := exec.Command(name, args...)
		var b bytes.Buffer
		cmd.Stdout = &b
		cmd.Stderr = &b
		err := cmd.Run()
		ch <- AsyncCommandResult{b.Bytes(), err}
	}(ch)
	return ch
}
