#! /usr/bin/env bash

set -euo pipefail

git config --global --add safe.directory '*'

cd /gocms/migrations
GOOSE_DRIVER="mysql" GOOSE_DBSTRING="root:my-secret-pw@tcp(mariadb:3306)/cms_db" goose up

cd /gocms
air

# root:my-secret-pw@tcp(localhost:3306)/cms_db