build: rencode_generated.go
	go build

tests: rencode_generated.go
	go test

clean:
	rm -f rencode_generated.go

rencode_generated.go:
	@rm -f rencode_generated.go
	go run generate.go > rencode_generated.go.tmp
	mv rencode_generated.go.tmp rencode_generated.go

.PHONY: build tests clean
