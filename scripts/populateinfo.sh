#!/bin/sh

FILE_LOCATION="./internal/static/embedded"

VERSION=$(git describe --tags --abbrev=0)
DATE="$(TZ='UTC' date) UTC"

[ "$VERSION" == "" ] && {
    COMMIT=$(git rev-parse HEAD)
    VERSION="c${COMMIT:0:8}"
}

printf "$VERSION" > "$FILE_LOCATION/version.txt"
printf "$DATE" > "$FILE_LOCATION/builddate.txt"