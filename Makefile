.DEFAULT: tests
.PHONY: tests

clean:
	rm -f rencode_generated.go

tests: rencode_generated.go
	go test

rencode_generated.go:
	@rm -f rencode_generated.go
	go generate > rencode_generated.go.tmp
	mv rencode_generated.go.tmp rencode_generated.go
