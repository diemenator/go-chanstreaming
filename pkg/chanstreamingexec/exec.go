package chanstreamingexec

import (
	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
)

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

			if n > 0 {
				bufCpy := make([]byte, n)
				copy(bufCpy, buf[:n])
				procOut := ProcOut{
					MessageType: messageType,
					DataBytes:   bufCpy,
					Time:        time.Now(),
				}
				out <- procOut
			}
		}
	}()
	return out
}

// FromCmdStdOut invokes cmd.StdoutPipe() to produce a readonly channel of captured standard output binary stream
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

// NewStdinPipeSink creates a function that takes a channel of ProcIn and communicates it to the running `cmd`
func NewStdinPipeSink(cmd *exec.Cmd, onWriteError WriteErrorCallback) (func(src <-chan ProcIn), error) {
	stdInPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
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
					onWriteError(nil)
				}
			}
		}
	}

	result := func(src <-chan ProcIn) {
		go func() {
			for ins := range src {
				processItem(ins)
			}
			closeErr := stdInPipe.Close()
			if closeErr != nil {
				onWriteError(closeErr)
			}
		}()
	}

	return result, nil
}

// StartCommand executes the passed in command and returns a channel of ProcOut, feeding the command's stdin with the provided source channel in background.
//
// The command's stdout and stderr are captured and returned in the resulting channel and the exit code as well as errors that occur.
//
// The command is started in the background and the function returns the channel immediately.
//
// The `writeErrorCallback` is called when an error occurs while writing to the command's stdin.
//
// The `writeErrorCallback` Used to introduce logging with backoff or panicking behavior.
//
// The resulting channel would typically be consisting of interleaved StdOut, StdErr and finally ExitCode messages.
//
// Any process startup or communications errors are returned with IOError messages as well.
func StartCommand(cmd *exec.Cmd, writeErrorCallback WriteErrorCallback, src <-chan ProcIn) <-chan ProcOut {
	cmdSink, sinkErr := NewStdinPipeSink(cmd, writeErrorCallback)
	if sinkErr != nil {
		return NewIOErrorChan(sinkErr)
	}
	ios := ch.Merged(FromCmdStdErr(cmd), FromCmdStdOut(cmd))
	startErr := cmd.Start()
	if startErr != nil {
		return ch.Concat(NewIOErrorChan(startErr), ios)
	}

	cmdSink(src)
	return ch.Concat(ios, FromProcAwait(cmd))
}

func NewProcSignal(sig os.Signal) ProcIn {
	return ProcIn{
		MessageType: Signal,
		Signal:      sig,
	}
}

func NewProcIn(b []byte, messageType ProcInType) ProcIn {
	return ProcIn{
		MessageType: messageType,
		DataBytes:   b,
	}
}

func NewProcInStr(str string, messageType ProcInType) ProcIn {
	return NewProcIn([]byte(str), messageType)
}

func NewProcStdinStr(str string) ProcIn {
	return NewProcInStr(str, StdIn)
}
