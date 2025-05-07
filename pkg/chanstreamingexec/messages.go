package chanstreamingexec

import (
	"os"
	"os/exec"
	"time"
)

type ProcOutType int

const (
	IOError ProcOutType = iota
	StdOut
	StdErr
	ExitCode
)

type ProcInType int

const (
	StdIn ProcInType = iota
	Signal
)

// ProcOut is a struct that represents the captured output of a process.
type ProcOut struct {
	MessageType ProcOutType
	Origin      *exec.Cmd
	Time        time.Time
	Error       error
	DataBytes   []byte
	ExitCode    int
}

// ProcIn is a struct that represents the input to a process.
type ProcIn struct {
	MessageType ProcInType
	DataBytes   []byte
	Signal      os.Signal
}
