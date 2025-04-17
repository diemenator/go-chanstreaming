# <-chan streaming exec

This package adds `StartCommand(*exec.Cmd, ...)` command-based channel transform function.

It accepts a channel of mix of raw IO bytes (to be written to stdin) and os.Signals to 
be communicated to the running command process.

