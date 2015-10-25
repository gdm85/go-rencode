.DEFAULT := tests
.PHONY := tests pregen

tests: pregen
	go test


pregen:
	rm -f rencode_generated.go
	go generate
