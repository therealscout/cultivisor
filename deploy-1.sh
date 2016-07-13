#!/bin/bash

clear
NODE="$1"
DIR=`echo ${PWD##*/}`

INCLUDE="${DIR} assets/ templates/"

if [ "$NODE" == "" ]; then
    echo "No NODE Specified"
    exit 1
fi

if [ -f "${DIR}.tar" ]; then
    echo "Removing Old Tar ${DIR}.tar"
    rm $DIR.tar
fi

echo "Remving Old Binary ${DIR}"
go clean

echo "Building ${DIR}"
go build

if [ ! -f $DIR ]; then
    echo "Build ${DIR} Failed"
    exit 1
fi

echo "Creating tar ${DIR}.tar"
tar cf $DIR.tar $INCLUDE

if [ ! -f "${DIR}.tar" ]; then
    echo "Creating ${DIR}.tar Failed"
    exit 1
fi

echo "SCP to NODE ${NODE}"
scp $DIR.tar node$NODE@node$NODE.cagnosolutions.com:/home/node$NODE
