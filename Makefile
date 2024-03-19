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
	GIN_MODE=release DATABASE_HOST=localhost DATABASE_PORT=3306 DATABASE_USER=root DATABASE_PASSWORD=my-secret-pw DATABASE_NAME=cms_db IMAGE_DIRECTORY="/home/enrique/Development/Go/gocms/media" CONFIG_FILE_PATH="settings/gocms_config.toml" $(OUT)/$(ADMINAPPNAME)

test:
	$(GOTEST) -v ./...

install-tools:
	go install github.com/a-h/templ/cmd/templ@v0.2.476
	go install github.com/pressly/goose/v3/cmd/goose@v3.18.0

run-containers:
	docker build -t emarifer/gocms:0.1  docker/
	docker compose up

.PHONY: clean
