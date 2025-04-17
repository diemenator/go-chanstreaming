package chanstreamingexec_test

import (
	"errors"
	"fmt"
	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	chexec "github.com/diemenator/go-chanstreaming/pkg/chanstreamingexec"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func logError(t *testing.T) func(err error) {
	return func(err error) { t.Error(err) }
}

func TestEcho(t *testing.T) {

	echoCommand := chexec.NewShellCommand("echo hello world")
	procIo := chexec.StartCommand(echoCommand, logError(t), ch.Empty[chexec.ProcIn]())
	slice := ch.ToSlice(procIo)
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

const sampleInput string = `hello world 1
hello world 2
	asd
hello world 3
hello world 4
hello world 4.5
hello world 5
`

const sampleOutput string = `You said: hello world 1
You said: hello world 2
You said: 	asd
You said: hello world 3
You said: hello world 4
You said: hello world 4.5
You said: hello world 5`

const sampleInputLineLength int = 7

func getRepeaterExecutable() string {
	echoBackScript := "echo-back"
	if runtime.GOOS == "windows" {
		echoBackScript = `echo-back.cmd`
	} else {
		echoBackScript = `echo-back.sh`
	}
	echoBackScript, err := filepath.Abs(echoBackScript)
	if err != nil {
		message := fmt.Sprint("failed to find absolute path for the test stdin 'You said:<>' repeater", echoBackScript, err)
		ourErr := errors.New(message)
		panic(ourErr)
	}
	return echoBackScript
}

func newTestStdIn() <-chan chexec.ProcIn {
	lines := strings.Split(sampleInput, "\n")
	linesChan := ch.FromSlice(lines)
	stdIn := ch.Mapped(func(x string) chexec.ProcIn {
		return chexec.NewProcStdinStr(x + "\n")
	})(linesChan)
	return stdIn
}

func TestEchoStdIn(t *testing.T) {
	stdIn := newTestStdIn()
	stdIn = ch.Throttle[chexec.ProcIn](time.Second)(stdIn)
	startTime := time.Now()
	procIo := chexec.StartCommand(exec.Command(getRepeaterExecutable()), logError(t), stdIn)

	capturedProcOutputs := ch.ToSlice(procIo)

	elapsed := time.Since(startTime)
	t.Log("Elapsed time:", elapsed)
	if elapsed < (time.Second * time.Duration(sampleInputLineLength-1)) {
		t.Error("Elapsed time is less than expected:", elapsed)
	}

	capturedStdStr := ""
	for _, v := range capturedProcOutputs {
		if v.MessageType == chexec.StdOut || v.MessageType == chexec.StdErr {
			capturedStdStr = capturedStdStr + string(v.DataBytes)
		}
	}
	capturedStdStr = strings.TrimSpace(capturedStdStr)
	capturedStdStr = strings.ReplaceAll(capturedStdStr, "\r", "")

	if strings.Compare(sampleOutput, capturedStdStr) == 0 {
		t.Log("Captured:", capturedStdStr)
	} else {
		t.Error("Expected:\n", sampleOutput, "\nGot:\n", capturedStdStr)
	}
}
