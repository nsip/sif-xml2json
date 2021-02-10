#!/bin/bash
set -e

R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
B=`tput setaf 4`
W=`tput sgr0`

printf "\n"

ip="localhost:1324/"          ###
base=$ip"sif-xml2json/v0.1.5" ###

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

sv=3.4.8.draft

SIFXFile=./data/examples/StudentPersonals.xml ###
title='Convert Test @ '$SIFXFile
url=$base"?sv=$sv&wrap"   ###
file="@"$SIFXFile
scode=`curl -X POST $url --data-binary $file -w "%{http_code}" -s -o /dev/null`
if [ $scode -ne 200 ]; then
    echo "${R}${title}${W}"
    exit 1
else
    echo "${G}${title}${W}"
fi

jsonname=`basename $SIFXFile .xml`"@$sv".json
outdir=./data/output/
mkdir -p $outdir
outfile=$outdir"$jsonname"
echo "curl -X POST $url --data-binary $file"
curl -X POST $url --data-binary $file > $outfile
cat $outfile
printf "\n"

echo "${G}Done${W}"

####################################################

# title='SIF-XML2JSON all API Paths'
# url=$ip
# scode=`curl --write-out "%{http_code}" --silent --output /dev/null $url`
# if [ $scode -ne 200 ]; then
#     echo "${R}Error getting root information from ${ip} - ${title}${W}"
#     exit 1
# else
#     echo "${G}Server OK: ${title}${W}"
# fi
# echo "# Headers: curl $url"
# curl -i $url
# printf "\n"

# # exit 0

# sv=3.4.7

# SIFDir=./data/examples/$sv/*
# for f in $SIFDir
# do
#     title='Convert Test @ '$f
#     url=$base"?sv=$sv"    ###
#     file="@"$f
#     scode=`curl -X POST $url --data-binary $file -w "%{http_code}" -s -o /dev/null`
#     if [ $scode -ne 200 ]; then
#         echo "${R}${title}${W}"
#         exit 1
#     else
#         echo "${G}${title}${W}"
#     fi

#     jsonname=`basename $f .xml`.json
#     outdir=./data/output/$sv/
#     mkdir -p $outdir
#     outfile=$outdir"$jsonname"
#     echo "curl -X POST $url --data-binary $file"
#     curl -X POST $url --data-binary $file > $outfile
#     cat $outfile
#     printf "\n"
# done

# echo "${G}All Done${W}"
