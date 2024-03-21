#! /usr/bin/env bash

set -euo pipefail

git config --global --add safe.directory '*'

cd /gocms/migrations
GOOSE_DRIVER="mysql" GOOSE_DBSTRING="root:my-secret-pw@tcp(mariadb:3306)/cms_db" goose up

cd /gocms
air -c ./docker/.air.toml

# root:my-secret-pw@tcp(localhost:3306)/cms_db

# reason this line (git config --global --add safe.directory '*') is added:
# https://www.youtube.com/watch?v=41iONer9RxM&list=PLZ51_5WcvDvCzCB2nwm8AWodXoBbaO3Iw&index=10&t=4063s
# https://www.google.com/search?q=detect+dubious+ownership+in+repository+at+docker+container&oq=detect+dubious+ownership+in+repository+at+docker+co&aqs=chrome.2.69i57j33i160l2j33i21.9509j0j4&sourceid=chrome&ie=UTF-8#ip=1
# https://community.jenkins.io/t/detected-dubious-ownership-in-repository-with-jenkins-upgrade/6182
# https://docs.docker.com/compose/compose-file/compose-file-v3/#volumes

# ============================================================================
# Explanation from this moment "https://www.youtube.com/live/41iONer9RxM?si=7SA8DwA9yyu_BhJi&t=7134" until 5 minutes later ("https://www.youtube.com/live/41iONer9RxM?si=BD0TT3Rgqk_FJxtp&t=7349")

# SCRIPT EXAMPLE THAT FAILS LOUDLY:
# #!/bin/bash

# set -euxo pipefail

# firstName="Aaron"
# fullName="$firstname Maxwell"
# echo "$fullName"

# IF WE DO NOT ADD THE "set -euxo pipefail" LINE THE SCRIPT WILL FAIL WITHOUT US KNOWING WHY.
# WITH THE COMMAND "echo $?" WE CAN OBTAIN THE EXIT CODE:
# echo $? =>
# 1
# [What does echo $? do?] => https://unix.stackexchange.com/questions/501128/what-does-echo-do
# REFERENCES:
# https://www.google.com/search?q=set+-euo+pipefail&oq=set+-euo+pipefail&aqs=chrome..69i57.1859395j0j7&sourceid=chrome&ie=UTF-8
# https://gist.github.com/mohanpedala/1e2ff5661761d3abd0385e8223e16425?permalink_comment_id=3935570