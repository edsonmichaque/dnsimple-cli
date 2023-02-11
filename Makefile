.PHONY: build
build:
	go build -o ./bin/dnsimple cmd/dnsimple/main.go

.PHONY: test
test:
	go test -race ./...

.PHONY: dep
make dep:
	go mod download

.PHONY: release
release:
	goreleaser release --clean


.PHONY: install-addlicense
install-addlicense:
	go install github.com/google/addlicense@latest

.PHONY: copyright
copyright: install-addlicense
	addlicense -c 'Edson Michaque' -y 2023 -l apache -s  -ignore .github/** -ignore *.yml -ignore *.yaml .
