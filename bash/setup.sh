#!bin/bash

echo "Creating docker network: dev_network"
docker network create --driver bridge dev_network