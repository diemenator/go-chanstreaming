module github.com/diemenator/go-chanstreaming/pkg/chanstreamingtests

go 1.24.1

require (
	github.com/diemenator/go-chanstreaming/pkg/chanstreaming v0.0.0-20250429082921-41ba5d19b0cf
	github.com/diemenator/go-chanstreaming/pkg/chanstreamingexec v0.0.0-20250429082921-41ba5d19b0cf
	github.com/stretchr/testify v1.10.0
)

replace (
	github.com/diemenator/go-chanstreaming/pkg/chanstreaming => ../pkg/chanstreaming
	github.com/diemenator/go-chanstreaming/pkg/chanstreamingexec => ../pkg/chanstreamingexec
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
