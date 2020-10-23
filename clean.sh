#!/bin/bash

set -e
shopt -s extglob

ORIPATH=`pwd`

# cd ./Server/ && ./clean.sh && cd $ORIPATH && echo "Server clean"
cd ./SIFSpec/ && ./clean.sh && cd $ORIPATH && echo "SIFSpec clean"

rm -rf ./data/output/
rm -f ./*.json ./*.xml ./go.sum

# delete all binary files
find . -type f -executable -exec sh -c "file -i '{}' | grep -q 'x-executable; charset=binary'" \; -print | xargs rm -f
for f in $(find ./ -name '*.log' -or -name '*.doc'); do rm $f; done