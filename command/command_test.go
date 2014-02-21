package command

import (
	"bytes"
	"testing"
)

func TestSyncCommand(t *testing.T) {
	output, err := SyncCommand("echo", "hello")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(bytes.TrimSpace(output), []byte("hello")) {
		t.Fail()
	}
}

func TestAsyncCommand(t *testing.T) {
	ch := AsyncCommand("echo", "hello")
	r := <-ch
	if r.Error != nil {
		t.Error(r.Error)
	}
	if !bytes.Equal(bytes.TrimSpace(r.Output), []byte("hello")) {
		t.Fail()
	}
}
