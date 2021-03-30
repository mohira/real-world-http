FROM golang:1.16

WORKDIR /src

RUN apt update && \
    apt install -y strace less vim && \
    apt install -y man manpages-ja manpages-ja-dev && \
    apt install -y netcat lsof


