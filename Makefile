.PHONY: build
build:
	go build -o ./bin/dnsimple cmd/dnsimple/main.go

.PHONY: test
test:
	go test -race -v ./...

.PHONY: dep
dep:
	go mod download

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: release
release:
	goreleaser release --clean

.PHONY: addlicense
addlicense:
	go install github.com/google/addlicense@latest

.PHONY: copyright
copyright: install-addlicense
	addlicense -c 'Edson Michaque' -y 2023 -l apache -s  .

.PHONY: check-license
check-license: addlicense
	addlicense -check .
