all: deps vet staticcheck

deps:
	go mod download

deps-update:
	go get

vet: deps
	go vet ./...

staticcheck: deps
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

test: deps
	go test

