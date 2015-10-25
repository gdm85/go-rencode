.DEFAULT: tests
.PHONY: tests

tests: rencode_generated.go
	go test

rencode_generated.go:
	go generate > rencode_generated.go
