#!/bin/bash
set -e

R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
B=`tput setaf 4`
W=`tput sgr0`

printf "\n"

ip="192.168.31.100:1324/"      ###
base=$ip"sif-xml2json/v0.1.2/" ###

title='SIF-XML2JSON all API Paths'
url=$ip
scode=`curl --write-out "%{http_code}" --silent --output /dev/null $url`
if [ $scode -ne 200 ]; then
    echo "${R}${title}${W}"
    exit 1
else
    echo "${G}${title}${W}"
fi
echo "curl $url"
curl -i $url
printf "\n"

# exit 0

sv=3.4.6

SIFXFile=./data/examples/siftest346.xml
title='Convert Test @ '$SIFXFile
url=$base"convert?sv=$sv&wrap"   ###
file="@"$SIFXFile
scode=`curl -X POST $url --data-binary $file -w "%{http_code}" -s -o /dev/null`
if [ $scode -ne 200 ]; then
    echo "${R}${title}${W}"
    exit 1
else
    echo "${G}${title}${W}"
fi

jsonname=`basename $SIFXFile .xml`.json
outdir=./data/output/
mkdir -p $outdir
outfile=$outdir"$jsonname"
echo "curl -X POST $url --data-binary $file"
curl -X POST $url --data-binary $file > $outfile
cat $outfile
printf "\n"

echo "${G}Done${W}"
