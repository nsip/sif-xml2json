#!/bin/bash
set -e

R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
B=`tput setaf 4`
W=`tput sgr0`

printf "\n"

ip="192.168.31.168:1324/"      ### 
base=$ip"sif-xml2json/v0.1.0/" ###

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

sv=3.4.7

SIFDir=./data/examples/$sv/*
for f in $SIFDir
do    
    title='Convert Test @ '$f
    url=$base"convert?sv=$sv"    ###
    file="@"$f
    scode=`curl -X POST $url -d $file -w "%{http_code}" -s -o /dev/null`
    if [ $scode -ne 200 ]; then
        echo "${R}${title}${W}"
        exit 1
    else
        echo "${G}${title}${W}"
    fi

    jsonname=`basename $f .xml`.json
    outdir=./data/output/$sv/
    mkdir -p $outdir
    outfile=$outdir"$jsonname"
    echo "curl -X POST $url -d $file"
    curl -X POST $url -d $file > $outfile
    cat $outfile
    printf "\n"
done

echo "${G}All Done${W}"