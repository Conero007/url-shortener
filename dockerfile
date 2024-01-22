FROM golang:1.21

RUN apt-get update && apt-get install -y wget vim telnet git && apt-get update

RUN apt-get update && apt-get install pkg-config libssl-dev net-tools -y && apt-get update

RUN git config --global --add safe.directory /var/www/url_shortener

RUN mkdir -p /var/www

WORKDIR /var/www