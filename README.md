# <-chan streaming

## TLDR

- install with `go get github.com/diemenator/go-chanstreaming/pkg/chanstreaming@latest`  

- read docs at: https://pkg.go.dev/github.com/diemenator/go-chanstreaming/pkg/chanstreaming#pkg-functions

# Packages

- #### [chanstreaming](pkg/chanstreaming/README.md)  
- #### [chanstreamingexec](pkg/chanstreamingexec/README.md)  

# See also:

https://github.com/reugn/go-streams

https://github.com/golang-design/go2generics/tree/main/chans  

# todo

- [ ] rebrand to `rochannels` for std lib

- [ ] reconsider or move into separate module some of controversial pieces of API like killswitches and 'async-like' names

- [ ] showcas


### go.work

```work
go 1.24.1

use (
	./pkg/chanstreaming
	./pkg/chanstreamingexec
	./tests
)
```
