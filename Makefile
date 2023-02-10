.PHONY: build
build:
	go build -o ./bin/dnsimple cmd/dnsimple/main.go

.PHONY: test
test:
	go test -race ./...

.PHONY: dep
make dep:
	go mod download
