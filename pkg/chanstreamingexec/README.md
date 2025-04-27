# <-chan streaming exec

This package adds single `StartCommand(*exec.Cmd, ...)` command-based channel transform function.

It creates a flow function - accepting a readonly channel of ProcIn messages - of mix of raw IO bytes (to be written to stdin) and os.Signals - and maps it to ProcOut - a typed mix out stderr, stdout and exitcode messages emitted as they come from the process underneath.

That way simple stdio transport-enabled app like (./echo-back.sh)[./echo-back.sh] (or your favourite mcp app) can be wrapped into coroutine that feeds ProcIn channel and a single StartCommand call.

See demo usage in tests (./chanstreamingexec/exec_test.go#L116-L120)[./chanstreamingexec/exec_test.go#L116-L120].



