#!/bin/sh

version=1.0.0
echo "qdmovie restart version:" $version

sh ./stopdog.sh
sh ./startdog.sh
