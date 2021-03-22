#!/bin/bash
# rm -f ./go.sum
# go get -u ./...

ORIPATH=`pwd`

cd ./config && ./build_d.sh && cd "$ORIPATH" 
echo "config prepared"

cd ./server && ./build_d.sh && cd "$ORIPATH" 
echo "server building done"
