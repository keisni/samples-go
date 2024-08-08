#!/bin/bash 
set -e
export GOOS=linux 
export GOARCH=amd64

WORKDIR=`pwd`
cd ./producer/starter
go build -buildvcs=false 


cd $WORKDIR
cd ./producer/worker
go build -buildvcs=false

cd $WORKDIR
cd ./consumer/starter
go build -buildvcs=false

cd $WORKDIR
cd ./consumer/worker
go build -buildvcs=false

