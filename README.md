# 0block

This service is responsible for fetching blocks and their transactions from blockchain and store them in a MongoDB instance, So that we can leverage that DB instance and use it for explorer or any other FE application which requires block and transactions data.

## Table of Contents

- [Setup](#setup)
- [Buildding and starting the node](#building-and-starting-the-node)
- [Point to another blockchain](#point-to-another-blockchain)
- [Exposed APIs](#exposed-apis)
- [Config changes](#config-changes)
- [Cleanup](#cleanup)
- [Network issue](#network-issue)

## Setup

Clone the repo and run the following command inside the cloned directory

```
$ ./docker.local/bin/init.sh
```

## Building and Starting the Node

If there is new code, do a git pull and run the following command

```
$ ./docker.local/bin/build.sh
```
You may need to run the following command occasionaly to make sure your device is in sync with ntp

```
$ ./docker.local/bin/sync_clock.sh
```

Go to the bin directory (cd docker.local/bin) and run the container using

```
$ ./start.sh
```

## Point to another blockchain

You can point the server to any instance of 0chain blockchain you like, Just go to config (docker.local/config) and update the 0block.yaml.

```
dns_url: http://198.18.0.98:9091
```

We use 0dns to connect to the network instead of giving network details directly, It will fetch the network details automatically from the 0dns's network API.

There are other configurable properties as well which you can update as per the requirement.

## Exposed APIs

### Logs APIs

```
{BASEURL}/logs
{BASEURL}/mem_logs
```

These APIs can be used to view logs or memory logs of the running 0block instance.
To view logs in detail use `{BASEURL}/logs?detail=3`, Same for `mem_logs`.

## Config Changes

You can do other config changes as well in 0block.yaml file itself, Like

- Mongo DB connection URL, DB name and pool size

```
mongo:
  url: mongodb://mongodb:27017
  db_name: block-recorder
  pool_size: 2
```

- Server port

```
port: 9091
```

- Logging info level

```
logging:
  level: "info"
  console: false # printing log to console is only supported in development mode
```

- Worker related config

```
worker:
  round_fetch_delay: 250 # in milliseconds, Wait before retrying for the same round
  round_fetch_retires: 15 # Retries to fetch a round (if failed)
```

## Cleanup

Get rid of old data when the blockchain is restarted or if you point to a different network:

```
$ ./docker.local/bin/clean.sh
```

## Network issue

If there is no test network, run the following command

```
docker network create --driver=bridge --subnet=198.18.0.0/15 --gateway=198.18.0.255 testnet0
```
