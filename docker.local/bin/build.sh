#!/bin/sh

docker-compose -p 0block -f docker.local/docker-compose.yml build --force-rm
docker.local/bin/sync_clock.sh
