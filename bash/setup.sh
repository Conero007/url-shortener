#!bin/bash

echo "Creating docker network: dev_network"
docker network create --driver bridge dev_network

echo "Creating .env file"
sh -c "touch .env && cat .example.env > .env"