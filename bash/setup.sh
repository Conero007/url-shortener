#!bin/bash

echo "Creating docker network: dev_network"
docker network create --driver bridge dev_network

echo "Creating .env file"
sh -c "touch .env && cat .example.env > .env"

echo "Updating /etc/hosts..."
ip_address = "127.0.0.1"
host_name="url.shortener.local"

matches_in_hosts = "$(grep -n "$hostname" /etc/hosts | cut -f1 -d:)"
host_entry = "${ip_address} ${host_name}"

if [ -n "$matches_in_hosts" ]
then
    echo "Host entry for [$host_name] alredy exists."
else
    echo "Adding new hosts entry: $host_name."
    echo "$host_entry" | sudo tee -a /etc/hosts > /dev/null
fi