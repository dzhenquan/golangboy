#!/bin/sh

version=1.0.0
echo "qdmovie start version:" $version

LOGDIR=/root/Golang/src/github.com/dzhenquan/golangboy

SERVICE1=golangboy
SERV1_DIR=/root/Golang/src/github.com/dzhenquan/golangboy

if [ ! -d $LOGDIR ]; then
    mkdir $LOGDIR
fi

while [ 1 ]
do
    progname=`ps -A| grep $SERVICE1 | grep -v grep | awk '{print $4}'`
    if [ -z $progname ]; then
       echo "$SERVICE1 process isn't running!!! restart it now."

       cd $SERV1_DIR
      sleep 4
       ./$SERVICE1 &
      sleep 1
    fi

    sleep 5
done

