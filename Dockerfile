FROM golang:1.17-bullseye

RUN apt-get update && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

WORKDIR /go/src/github.com/Fukkatsuso/twego

COPY go.* ./
RUN go mod download
