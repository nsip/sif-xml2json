#!/bin/bash
set -e

mkdir -p ./app
cp ./server/build/linux64/* ./app/
echo "all built"