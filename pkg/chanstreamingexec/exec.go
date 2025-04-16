package chanstreamingexec

import (
	"errors"
	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type ProcOutType int // 0 - InIO, 1 - OutIO, 2 - InSignal, 3 - InData, 4 - ProcIO
const (
	IOError ProcOutType = iota
	StdOut
	StdErr
	ExitCode
)

type ProcOut struct {
	Origin      *exec.Cmd
	Time        time.Time
	MessageType ProcOutType
	Error       error
	DataBytes   []byte
	ExitCode    int
}

type ProcInType int // 0 - InIO, 1 - OutIO, 2 - InSignal, 3 - InData, 4 - ProcIO
const (
	StdIn ProcInType = iota
	Signal
)

type ProcIn struct {
	MessageType ProcInType
	DataBytes   []byte
	Signal      os.Signal
}

func NewShellCommand(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/C", command)
	} else {
		return exec.Command("sh", "-c", command)
	}
}

func FromIoReadCloser(reader *io.ReadCloser, messageType ProcOutType) <-chan ProcOut {
	if reader == nil {
		return nil
	}
	theReader := *reader
	out := make(chan ProcOut)

	// Ensure the reader is not nil
	go func() {
		defer func() {
			close(out)
		}()
		buf := make([]byte, 1024)
		for {
			n, err := theReader.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				out <- ProcOut{
					MessageType: IOError,
					Error:       err,
					Time:        time.Now(),
				}
				return
			}
			out <- ProcOut{
				MessageType: messageType,
				DataBytes:   buf[:n],
				Time:        time.Now(),
			}
		}
	}()
	return out
}

// FromCmdStdOut invokes cmd.StdoutPipe() to produce a readonly channel of captured standard error binary stream
func FromCmdStdOut(cmd *exec.Cmd) <-chan ProcOut {
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return NewIOErrorChan(err)
	} else {
		return FromIoReadCloser(&reader, StdOut)
	}
}

// FromCmdStdErr invokes cmd.StderrPipe() to produce a readonly channel of captured standard error binary stream
func FromCmdStdErr(cmd *exec.Cmd) <-chan ProcOut {
	reader, err := cmd.StderrPipe()
	if err != nil {
		return NewIOErrorChan(err)
	} else {
		return FromIoReadCloser(&reader, StdErr)
	}
}

// FromProcAwait asynchronously invokes `cmd.Wait()`
// and emits ExitCode
// or IOError with received io.ExitError
// in resulting channel
func FromProcAwait(cmd *exec.Cmd) <-chan ProcOut {
	out := make(chan ProcOut)
	go func() {
		defer func() {
			close(out)
		}()

		if cmd == nil {
			return
		}

		if err := cmd.Wait(); err != nil {
			out <- ProcOut{
				MessageType: IOError,
				Origin:      cmd,
				Error:       err,
				Time:        time.Now(),
			}
			return
		}

		out <- ProcOut{
			MessageType: ExitCode,
			Origin:      cmd,
			ExitCode:    cmd.ProcessState.ExitCode(),
			Time:        time.Now(),
		}
	}()
	return out
}

func NewIOError(err error) ProcOut {
	return ProcOut{MessageType: IOError, Error: err}
}

func NewIOErrorSlice(err error) []ProcOut {
	return []ProcOut{NewIOError(err)}
}

func NewIOErrorChan(err error) <-chan ProcOut {
	return ch.FromSlice[ProcOut](NewIOErrorSlice(err))
}

type WriteErrorCallback func(error)

func IgnoreError(err error) {
	return
}

func ToCmdProc(cmd *exec.Cmd, src <-chan ProcIn, onWriteError WriteErrorCallback) error {
	stdInPipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	processItem := func(x ProcIn) {
		if x.MessageType == Signal {
			sigError := cmd.Process.Signal(x.Signal)
			if sigError != nil {
				onWriteError(sigError)
			}
		} else if x.MessageType == StdIn {
			toWrite := x.DataBytes
			// write
			for len(toWrite) > 0 {
				written, writeErr := stdInPipe.Write(toWrite)
				if writeErr != nil {
					onWriteError(writeErr)
				}
				if written > 0 {
					toWrite = toWrite[written:]
				} else {
					onWriteError(errors.New("idle"))
				}
			}
		}
	}

	go func() {
		for ins := range src {
			processItem(ins)
		}
	}()

	return nil
}

func Launch(cmd *exec.Cmd, src <-chan ProcIn, writeError WriteErrorCallback) <-chan ProcOut {
	writeErr := ToCmdProc(cmd, src, writeError)
	if writeErr != nil {
		writeError(writeErr)
	}
	ios := ch.Merged(FromCmdStdErr(cmd), FromCmdStdOut(cmd))
	startErr := cmd.Start()
	if startErr != nil {
		return ch.Concat(NewIOErrorChan(startErr), ios)
	}

	return ch.Concat(ios, FromProcAwait(cmd))
}
