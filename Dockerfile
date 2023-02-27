#################################
# Dockerfile to build & test gwm
#################################
FROM golang:1.17-alpine

MAINTAINER e. alvarez <55966724+ealvar3z@users.noreply.github.com>

RUN apk update && apk add xorg-server-dev xcb-util-dev

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go build -o /usr/bin/gwm ./cmd

CMD ["/usr/bin/gwm"]

