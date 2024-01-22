FROM golang:1.21

RUN apt-get update --fix-missing
RUN apt-get install -y wget vim telnet git pkg-config libssl-dev net-tools make
RUN apt-get update --fix-missing

RUN git config --global --add safe.directory /var/www/url_shortener

RUN mkdir -p /var/www

COPY . /var/www/url_shortener

WORKDIR /var/www/url_shortener

RUN make build
