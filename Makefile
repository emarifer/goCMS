#!make
include .dev.env

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
all: build test

# Migrations for DB
prepare_env:
	cp -r migrations tests/system_tests/helpers/

# ==== Development ====

install-tools:
	go install github.com/a-h/templ/cmd/templ@v0.2.476
	go install github.com/pressly/goose/v3/cmd/goose@v3.18.0
	go install github.com/cosmtrek/air@v1.49.0

start-devdb:
	docker compose -f docker/mariadb.yml up -d

run-migrations:
	GOOSE_DRIVER=mysql GOOSE_DBSTRING="${MARIADB_USER}:${MARIADB_ROOT_PASSWORD}@tcp(${MARIADB_ADDRESS}:${MARIADB_PORT})/${MARIADB_DATABASE}" goose -dir ./migrations up

build:
	$(TEMPLCMD) generate
	$(GOBUILD) -ldflags="-s -w" -v -o $(OUT)/$(APPNAME) $(SRC)/*.go
	$(GOBUILD) -ldflags="-s -w" -v -o $(OUT)/$(ADMINAPPNAME) $(ADMINSRC)/*.go

run:
#	$(GOBUILD) -o $(OUT)/$(APPNAME) $(SRC)/*.go
#	$(OUT)/$(APPNAME)
	DATABASE_HOST=localhost DATABASE_PORT=3306 DATABASE_USER=root DATABASE_PASSWORD=my-secret-pw DATABASE_NAME=cms_db IMAGE_DIRECTORY="./media" CONFIG_FILE_PATH="settings/gocms_config.toml" $(OUT)/$(ADMINAPPNAME)

# ==== Docker Containers ====

run-containers:
	docker build -t emarifer/gocms:0.1  docker/
	docker compose -f docker/docker-compose.yml up

start-admin-container:
	$(OUT)/$(ADMINAPPNAME) --config docker/gocms_config.toml

# Testing
test: prepare_env
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(OUT)

.PHONY: all build test clean

# Why does make think the target is up to date?. SEE:
# https://stackoverflow.com/questions/3931741/why-does-make-think-the-target-is-up-to-date
