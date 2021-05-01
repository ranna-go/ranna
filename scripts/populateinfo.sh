#!/bin/sh

FILE_LOCATION="./internal/static/embedded"

VERSION=$(git describe --tags --abbrev=0)
DATE="$(TZ='UTC' date) UTC"

[ "$VERSION" == "" ] && {
    COMMIT=$(git rev-parse HEAD)
    VERSION="c${COMMIT:0:8}"
}

echo "$FILE_LOCATION/version.txt"
printf "$VERSION" | tee "$FILE_LOCATION/version.txt"
echo ""
echo "CHECK: $(cat $FILE_LOCATION/version.txt)"

echo ""
echo "$FILE_LOCATION/builddate.txt"
printf "$DATE" | tee "$FILE_LOCATION/builddate.txt"
echo ""
echo "CHECK: $(cat $FILE_LOCATION/builddate.txt)"
