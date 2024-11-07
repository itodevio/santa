define RELEASE_BODY
See [README.md](https://github.com/itodevio/santa/blob/${VERSION}/README.md) for detailed instructions on how to use it and roadmap.

## Installation

#### Linux / macOS

For Linux or macOS users, you can quickly install and run Santa CLI using the following command:

```bash
curl -L https://github.com/itodevio/santa/releases/download/${VERSION}/`uname -s`/`uname -m`/santa -o /usr/local/bin/santa
chmod +x /usr/local/bin/santa
```

#### Windows

For Windows users, you can download the executable depending on your architecture:

| OS      | Architecture   | Download Link
|---------|----------------|----------------------------------
| Windows | amd64 (64-bit) | [Download](https://github.com/itodevio/santa/releases/download/${VERSION}/santa-amd64-Windows.exe
| Windows | arm64 (64-bit) | [Download](https://github.com/itodevio/santa/releases/download/${VERSION}/santa-arm64-Windows.exe
| Windows | 386 (32-bit)   | [Download](https://github.com/itodevio/santa/releases/download/${VERSION}/santa-i386-Windows.exe
---

*If unsure, you can usually download the `amd64` version, as it is widely compatible with most systems.
**Then, add the executable to your PATH.**

#### Build it from source

You can also build it from source by cloning the repository and running the following command:
```bash
make build
```

##### Happy coding and enjoy the Advent of Code! 🎅🎄
endef
export RELEASE_BODY

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

release-body:
	echo "$$RELEASE_BODY" > RELEASE.md
