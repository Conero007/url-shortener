FROM golang:1.21-alpine

RUN apk update && \
    apk add --no-cache git make && \
    rm -rf /var/cache/apk/*

RUN git config --global --add safe.directory /var/www/url_shortener

RUN mkdir -p /var/www/url_shortener

COPY . /var/www/url_shortener

WORKDIR /var/www/url_shortener
