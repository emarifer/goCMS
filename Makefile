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
	$(GOBUILD) -ldflags="-s -w" -v -o $(OUT)/$(APPNAME) $(SRC)/*.go
	$(GOBUILD) -ldflags="-s -w" -v -o $(OUT)/$(ADMINAPPNAME) $(ADMINSRC)/*.go

clean:
	$(GOCLEAN)
	rm -rf $(OUT)

run:
#	$(GOBUILD) -o $(OUT)/$(APPNAME) $(SRC)/*.go
#	$(OUT)/$(APPNAME)
	GIN_MODE=release $(OUT)/$(ADMINAPPNAME) --config settings/gocms_config.toml

test:
	$(GOTEST) -v ./...

.PHONY: clean
