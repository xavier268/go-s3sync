#!/bin/bash

BASE=$(pwd)
BASE=${BASE/go-s3sync*/}"go-s3sync"
echo "Base folder detected as :" $BASE

######################################################
## We are doing all of the above to ensure scripts  ##
## can run from anywhere in the main folder ;-)     ##
######################################################

MYBIN=$BASE/bin/
echo "Looking for $MYBIN"
if [ ! -d $MYBIN ] 
then
    echo "Creating $MYBIN"
    mkdir $MYBIN
fi

go build -o $MYBIN $BASE/cmd/...
ls -la $MYBIN
