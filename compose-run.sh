#! /usr/bin/env bash

set -euo pipefail

git config --global --add safe.directory '*'

cd /gocms/migrations
GOOSE_DRIVER="mysql" GOOSE_DBSTRING="root:my-secret-pw@tcp(mariadb:3306)/cms_db" goose up

cd /gocms
air

# root:my-secret-pw@tcp(localhost:3306)/cms_db

# reason this line (git config --global --add safe.directory '*') is added:
# https://www.youtube.com/watch?v=41iONer9RxM&list=PLZ51_5WcvDvCzCB2nwm8AWodXoBbaO3Iw&index=10&t=4063s
# https://www.google.com/search?q=detect+dubious+ownership+in+repository+at+docker+container&oq=detect+dubious+ownership+in+repository+at+docker+co&aqs=chrome.2.69i57j33i160l2j33i21.9509j0j4&sourceid=chrome&ie=UTF-8#ip=1
# https://community.jenkins.io/t/detected-dubious-ownership-in-repository-with-jenkins-upgrade/6182
# https://docs.docker.com/compose/compose-file/compose-file-v3/#volumes