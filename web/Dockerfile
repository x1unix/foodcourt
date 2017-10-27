# Voracity go-server docker file, version 1.0.0
FROM golang:1.9.2
MAINTAINER Denys Sedchenko <dsedchenko@llnw.com>

RUN apt-get update
COPY ./hosts  /etc/hosts

ENTRYPOINT /bin/bash
