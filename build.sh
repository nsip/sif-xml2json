#!/bin/bash
rm -f ./go.sum
go get -u ./...

R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
W=`tput sgr0`

oripath=`pwd`

cd ./config && ./build.sh && cd $oripath && echo "${G}config prepared${W}"
cd ./server && ./build.sh && cd $oripath && echo "${G}server building done${W}"
