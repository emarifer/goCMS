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
ADMINAPPNAME=gocms-admin

# Directories
SRC=./cmd/gocms
ADMINSRC=./cmd/gocms_admin
OUT=./tmp

# Targets
all: clean build

build:
	$(TEMPLCMD) generate
	$(GOBUILD) -v -o $(OUT)/$(APPNAME) $(SRC)/*.go
	$(GOBUILD) -v -o $(OUT)/$(ADMINAPPNAME) $(ADMINSRC)/*.go

clean:
	$(GOCLEAN)
	rm -rf $(OUT)

run:
#	$(GOBUILD) -o $(OUT)/$(APPNAME) $(SRC)/*.go
#	$(OUT)/$(APPNAME)
	$(OUT)/$(ADMINAPPNAME)

test:
	$(GOTEST) -v ./...

.PHONY: clean
