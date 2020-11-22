#!/bin/bash
rm -f ./go.sum
go get -u ./...

oripath=`pwd`

cd ./config && ./build_d.sh && cd $oripath && echo "config prepared"
cd ./server && ./build_d.sh && cd $oripath && echo "server building done"
