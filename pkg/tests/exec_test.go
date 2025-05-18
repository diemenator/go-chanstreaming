package chanstreamingtests_test

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	chexec "github.com/diemenator/go-chanstreaming/pkg/chanstreamingexec"
	"github.com/stretchr/testify/assert"
)

type LoggedErrors []error

func (l *LoggedErrors) logError(err error) {
	*l = append(*l, err)
}

func (l *LoggedErrors) complainIfAnyLoggedErrors(t *testing.T) {
	if len(*l) > 0 {
		t.Error("Error was logged:", *l)
	}
}

func TestEcho(t *testing.T) {
	echoCommand := chexec.NewShellCommand("echo hello world")
	loggedErrors := LoggedErrors{}
	procInput := ch.Empty[chexec.ProcIn]()

	procOutput := chexec.StartCommand(echoCommand, loggedErrors.logError, procInput)
	procOutSlice := ch.ToSlice(procOutput)
	procOutStr := collectProcOutSliceAsString(procOutSlice)
	t.Log(procOutStr)

	assert.Equal(t, 2, len(procOutSlice))
	stdoutMessage, exitMessage := procOutSlice[0], procOutSlice[1]
	if stdoutMessage.MessageType != chexec.StdOut {
		t.Error("wrong message type", stdoutMessage.MessageType)
	}
	stdoutMessageStr := strings.TrimSpace(string(stdoutMessage.DataBytes))
	if stdoutMessageStr != "hello world" {
		t.Error(stdoutMessageStr)
	}
	if exitMessage.MessageType != chexec.ExitCode {
		t.Error("wrong message type", exitMessage.MessageType)
	}
	if exitMessage.ExitCode != 0 {
		t.Error("wrong exit code", exitMessage.ExitCode)
	}
}

const sampleInput string = `hello world 1
hello world 2
	asd
hello world 3
hello world 4
hello world 4.5
hello world 5`

const sampleOutput string = `You said: hello world 1
You said: hello world 2
You said: 	asd
You said: hello world 3
You said: hello world 4
You said: hello world 4.5
You said: hello world 5`

const sampleInputLineLength int = 7

func getRepeaterExecutable() *exec.Cmd {
	scriptFileName := "echo-back"
	isWindows := runtime.GOOS == "windows"
	scriptFile := scriptFileName + ".sh"
	if isWindows {
		scriptFile = scriptFileName + ".cmd"
	}

	scriptPath, err := filepath.Abs(scriptFile)
	if err != nil {
		message := fmt.Sprint("failed to find absolute path for the test stdin 'You said:<>' repeater", scriptFileName, err)
		ourErr := errors.New(message)
		panic(ourErr)
	}

	if isWindows {
		return exec.Command(scriptPath)
	} else {
		return exec.Command("bash", scriptPath)
	}
}

func newSampleInputChannel() <-chan chexec.ProcIn {
	lines := strings.Split(sampleInput, "\n")
	linesChan := ch.FromSlice(lines)
	stdIn := ch.Mapped(func(x string) chexec.ProcIn {
		return chexec.NewProcStdinStr(x + "\n")
	})(linesChan)
	return stdIn
}

func collectProcOutSliceAsString(procIo []chexec.ProcOut) string {
	capturedStdStr := ""
	for _, v := range procIo {
		if v.MessageType == chexec.StdOut || v.MessageType == chexec.StdErr {
			capturedStdStr = capturedStdStr + string(v.DataBytes)
		}
	}
	capturedStdStr = strings.TrimSpace(capturedStdStr)
	capturedStdStr = strings.ReplaceAll(capturedStdStr, "\r", "")
	return capturedStdStr
}

func TestEchoStdIn(t *testing.T) {
	procInput := newSampleInputChannel()
	procInput = ch.Throttle[chexec.ProcIn](time.Second)(procInput)
	startTime := time.Now()
	loggedErrors := LoggedErrors{}
	procOutput := chexec.StartCommand(getRepeaterExecutable(), loggedErrors.logError, procInput)
	capturedProcOutputs := ch.ToSlice(procOutput)

	elapsed := time.Since(startTime)
	assert.GreaterOrEqual(t, elapsed, time.Second*time.Duration(max(0, sampleInputLineLength-1)))

	capturedStdStr := collectProcOutSliceAsString(capturedProcOutputs)
	loggedErrors.complainIfAnyLoggedErrors(t)
	assert.Equal(t, sampleOutput, capturedStdStr)
}

func newProcInputStreamWithSigkillInTheMiddle() <-chan chexec.ProcIn {
	linesChan1 := newSampleInputChannel()
	sigint := ch.FromSlice([]chexec.ProcIn{chexec.NewProcSignal(syscall.SIGKILL)})
	linesChan2 := newSampleInputChannel()
	return ch.Concat(linesChan1, sigint, linesChan2)
}

func TestSignal(t *testing.T) {
	procInput := newProcInputStreamWithSigkillInTheMiddle()
	procInput = ch.Throttle[chexec.ProcIn](time.Second)(procInput)
	loggedErrors := LoggedErrors{}

	procOutput := chexec.StartCommand(getRepeaterExecutable(), loggedErrors.logError, procInput)
	capturedProcOutputs := ch.ToSlice(procOutput)
	capturedStdStr := collectProcOutSliceAsString(capturedProcOutputs)
	loggedErrors.complainIfAnyLoggedErrors(t)

	assert.Equal(t, sampleOutput, capturedStdStr)
}
