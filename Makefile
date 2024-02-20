# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Application name
APPNAME=gocms

# Directories
SRC=./cmd/gocms
OUT=./tmp

# Targets
all: clean build

build:
	$(GOBUILD) -v -o $(OUT)/$(APPNAME) $(SRC)/*.go

clean:
	$(GOCLEAN)
	rm -f $(OUT)/$(APPNAME)

run:
	$(GOBUILD) -o $(OUT)/$(APPNAME) $(SRC)/*.go
	$(OUT)/$(APPNAME)

test:
	$(GOTEST) -v ./...

.PHONY: clean
