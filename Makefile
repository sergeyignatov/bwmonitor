COMMIT := $(shell git rev-parse HEAD)
all:
	go build -ldflags "-X main.revision=$(COMMIT)" -o bin/bwmonitor cmd/cmd.go
clean:
	rm -rf bin
