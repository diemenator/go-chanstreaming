package chanstreamingexec_test

import (
	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	chexec "github.com/diemenator/go-chanstreaming/pkg/chanstreamingexec"
	"strings"
	"testing"
)

func TestEcho(t *testing.T) {
	echoCommand := chexec.NewShellCommand("echo hello world")
	launched := chexec.Launch(echoCommand, ch.Empty[chexec.ProcIn](), chexec.IgnoreError)
	slice := ch.ToSlice(launched)
	t.Log(slice)

	if len(slice) == 2 {
		captured := strings.TrimSpace(string(slice[0].DataBytes))
		if captured != "hello world" {
			t.Error(captured)
		}
		exit := slice[1].ExitCode
		if exit != 0 {
			t.Error(exit)
		}
	} else {
		t.Error(len(slice))
	}
}
