# Go parameters
GOCMD=go
GOGET=$(GOCMD) get
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOPATH_=$(shell echo $(GOPATH))

# Other params
BINDATA_URL=github.com/jteeuwen/go-bindata/...
BINDATA=$(GOPATH_)/bin/go-bindata
 
# All are .PHONY for now because dependencyness is hard
.PHONY: default deps

deps:build
	$(GOGET) $(BINDATA_URL)
build:
	$(BINDATA) -o embed.go embed/
	$(GOBUILD) -o bin/cssmate
clean:
	$(GOCLEAN)
install:
	$(GOINSTALL)