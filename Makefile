GOFMT_FILES?=$$(find . -name '*.go' | findstr -v vendor)

fmt:
	gofmt -w $(GOFMT_FILES)

lint:
	golangci-lint run -v ./...

all: fmt lint