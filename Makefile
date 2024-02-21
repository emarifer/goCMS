# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Templ parameters
TEMPLCMD=templ

# Application name
APPNAME=gocms

# Directories
SRC=./cmd/gocms
OUT=./tmp

# Targets
all: clean build

build:
	$(TEMPLCMD) generate
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
