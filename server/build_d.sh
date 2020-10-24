#!/bin/bash
set -e

CGO_ENABLED=0 go run ./main.go ./var.go -- trial

rm -rf ./build

GOARCH=amd64
LDFLAGS="-s -w"
OUT=server

# For Docker, one build below for linux64 is enough.
OUTPATH=./build/linux64/
mkdir -p $OUTPATH
CGO_ENABLED=0 GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT $OUTPATH
cp ./config_rel.toml $OUTPATH'config.toml'
echo "server(linux64) built"

rm config_rel.toml
