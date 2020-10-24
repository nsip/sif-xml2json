#!/bin/bash
rm -f ./go.sum
go get -u ./...

oripath=`pwd`

cd ./sif-spec && ./build_d.sh && cd $oripath && echo "sif-spec ready"
cd ./config && ./build_d.sh && cd $oripath && echo "config prepared"
cd ./server && ./build_d.sh && cd $oripath && echo "server building done"
