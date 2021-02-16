## docker image prune
## docker rmi $(docker images -a -q)

# FROM alpine
# RUN mkdir /sif-xml2json
# COPY . / /sif-xml2json/
# WORKDIR /sif-xml2json/
# CMD ["./server"]

### ! run this Dockerfile 
### docker build --tag=sif-xml2json . 

### ! run this docker image
### docker run --name sif-xml2json --net host sif-xml2json:latest

### ! push image to docker hub
### docker tag IMAGE_ID dockerhub-user/sif-xml2json:latest
### docker login
### docker push dockerhub-user/sif-xml2json


###########################
# INSTRUCTIONS
############################
# BUILD
#	docker build --rm -t nsip/sif-xml2json:latest -t nsip/sif-xml2json:v0.1.0 .
# TEST: docker run -it -v $PWD/test/data:/data -v $PWD/test/config.json:/config.json nsip/sif-xml2json:develop .
# RUN: docker run -d nsip/sif-xml2json:develop
#
# PUSH
#	Public:
#		docker push nsip/sif-xml2json:v0.1.0
#		docker push nsip/sif-xml2json:latest
#
#	Private:
#		docker tag nsip/sif-xml2json:v0.1.0 the.hub.nsip.edu.au:3500/nsip/sif-xml2json:v0.1.0
#		docker tag nsip/sif-xml2json:latest the.hub.nsip.edu.au:3500/nsip/sif-xml2json:latest
#		docker push the.hub.nsip.edu.au:3500/nsip/sif-xml2json:v0.1.0
#		docker push the.hub.nsip.edu.au:3500/nsip/sif-xml2json:latest
#
###########################
# DOCUMENTATION
############################



# docker build --rm -t nsip/sif-xml2json:latest -t nsip/sif-xml2json:v0.1.0 .

###########################
# STEP 0 Get them certificates
############################
# (note, step 2 is using alpine now) 
# FROM alpine:latest as certs

############################
# STEP 1 build executable binary (go.mod version)
############################
FROM golang:1.15.8-alpine3.12 as builder
RUN apk add --no-cache ca-certificates
RUN apk update && apk add --no-cache git bash
RUN mkdir -p /sif-xml2json
COPY . / /sif-xml2json/
WORKDIR /sif-xml2json/
RUN ["/bin/bash", "-c", "./build_d.sh"]
RUN ["/bin/bash", "-c", "./release_d.sh"]

############################
# STEP 2 build a small image
############################
FROM alpine
COPY --from=builder /sif-xml2json/app/ /
# NOTE - make sure it is the last build that still copies the files
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /
CMD ["./server"]

# docker run --rm --mount type=bind,source=$(pwd)/config.toml,target=/config.toml -p 0.0.0.0:1324:1324 nsip/sif-xml2json