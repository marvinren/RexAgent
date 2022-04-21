.PHONY: all build clean run check cover lint docker help linuxbuild windowsbuild
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=bin/rexagent
MAIN_FILE=cmd/rexagent.go

all:test build

clean:
	$(GOCLEAN)
	rm -rf bin/*

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE) 

test:
	$(GOTEST) -v ./...

linuxbuild:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/rexagent_linux_amd64 $(MAIN_FILE)

windowsbuild:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/rexagent_windows.exe $(MAIN_FILE)
