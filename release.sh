#!/bin/bash
set -e

R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
W=`tput sgr0`

if [ $# -lt 2 ]; then
    echo "${Y}WARN:${W} input ${Y}Dest-OS-Type${W} [linux64 mac win64] and ${Y}Release Directory${W}"
    exit 1
fi

os=$1
dir=$2

if [ $os != 'linux64' ] && [ $os != 'mac' ] && [ $os != 'win64' ]; then
    echo "${Y}WARN:${W} input Dest-OS-Type [${Y}linux64 mac win64${W}]"
    exit 1
fi

mkdir -p $dir
cp ./server/build/$os/* $dir
echo "server($os) has been dumped into $dir"