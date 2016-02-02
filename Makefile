tests: rencode_generated.go
	go test

clean:
	rm -f rencode_generated.go

rencode_generated.go:
	@rm -f rencode_generated.go
	go generate > rencode_generated.go.tmp
	mv rencode_generated.go.tmp rencode_generated.go

.PHONY: tests
