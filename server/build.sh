#!/bin/bash
set -e

R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
W=`tput sgr0`

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
echo "${G}server(linux64) built${W}"

OUTPATH=./build/win64/
mkdir -p $OUTPATH
CGO_ENABLED=0 GOOS="windows" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT.exe
mv $OUT.exe $OUTPATH
cp ./config_rel.toml $OUTPATH'config.toml'
echo "${G}server(win64) built${W}"

OUTPATH=./build/mac/
mkdir -p $OUTPATH
CGO_ENABLED=0 GOOS="darwin" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT $OUTPATH
cp ./config_rel.toml $OUTPATH'config.toml'
echo "${G}server(mac) built${W}"

# GOARCH=arm
# OUTPATH=./build/linuxarm/
# mkdir -p $OUTPATH
# CGO_ENABLED=0 GOOS="linux" GOARCH="$GOARCH" GOARM=7 go build -ldflags="$LDFLAGS" -o $OUT
# mv $OUT $OUTPATH
# cp ./config_rel.toml $OUTPATH'config.toml'
# echo "${G}server(linuxArm) built${W}"

rm config_rel.toml