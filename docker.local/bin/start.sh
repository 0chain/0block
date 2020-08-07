#!/bin/sh
PWD=`pwd`

echo Starting 0block ...

docker-compose -p 0block -f ../docker-compose.yml up
