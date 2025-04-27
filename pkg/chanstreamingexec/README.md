# <-chan streaming exec

This package adds single `StartCommand(*exec.Cmd, ...)` command-based channel transform function and some channel types to model IO around it.

Stdio transport-enabled app like [./echo-back.sh](./echo-back.sh) (or your favourite mcp app) can be transparently invoked with a coroutine that feeds `<-chan ProcIn ` channel into `StartCommand`, resulting in `<-chan ProcOut` output.

See demo usage in tests [./exec_test.go#L116-L120](./exec_test.go#L116-L120).
