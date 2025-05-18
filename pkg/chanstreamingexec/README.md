# <-chan streaming exec

This package adds single `StartCommand(*exec.Cmd, ...)` command-based channel transform function and some channel types to model IO around it.

Stdio transport-enabled app like [tests/echo-back.sh](../tests/echo-back.sh) (or your favourite mcp app) can be transparently invoked with a coroutine that feeds `<-chan ProcIn ` channel into `StartCommand`, resulting in `<-chan ProcOut` output.

See demo usage in tests [tests/exec_test.go#L116-L120](../tests/exec_test.go#L116-L120).
