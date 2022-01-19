#!/bin/bash

OUTDIR="./docs/api"

# https://github.com/swaggo/swag
swag init \
    -g ./internal/api/v1/routes.go \
    -o $OUTDIR \
    --parseDependency --parseDepth 2

rm -rf $OUTDIR/docs.go

# https://github.com/syroegkin/swagger-markdown
swagger-markdown \
    -i $OUTDIR/swagger.json \
    -o $OUTDIR/restapi.md