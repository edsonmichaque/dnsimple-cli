COPYRIGHT_HOLDER := "Edson Michaque"
COPYRIGHT_YEARS := "2023"

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
copyright: addlicense
	addlicense -c ${COPYRIGHT_HOLDER} -y ${COPYRIGHT_YEARS} -l apache -s  .

.PHONY: check-license
check-license: addlicense
	addlicense -check .
