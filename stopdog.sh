#!/bin/sh

version=1.0.0
echo "qdmovie stop version:" $version

LOGDIR=/root/Golang/src/github.com/dzhenquan/golangboy

SERVICE1=golangboy
SERV1_DIR=/root/Golang/src/github.com/dzhenquan/golangboy

if [ ! -d $LOGDIR ]; then
    mkdir $LOGDIR
fi

progname=`ps -A| grep $SERVICE1 | grep -v grep | awk '{print $4}'`
if [ -z $progname ]; then
    echo "$SERVICE1 process isn't running!!!"
else
    echo "$SERVICE1 process is running!!! killed it now."
    killall $progname
fi

