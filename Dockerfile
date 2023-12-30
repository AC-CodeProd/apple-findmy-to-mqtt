FROM golang:1.20.4-alpine3.18 as base
ARG TZ
ARG UID
ARG GID
ARG USER=app

ENV USER=${USER}
ENV TZ=${TZ}
ENV UID=${UID}
ENV GID=${GID}
ENV GOCACHE /go/src/apple-findmy-to-mqtt/tmp/.cache
ENV GOLANGCI_LINT_CACHE /go/src/apple-findmy-to-mqtt/tmp/.cache

RUN addgroup -g $GID $USER && adduser -u $UID -G $USER -s /bin/sh -D $USER
RUN mkdir -p /go/src/apple-findmy-to-mqtt
WORKDIR /go/src/apple-findmy-to-mqtt
COPY . .
RUN go mod download

FROM base as development
RUN apk --update add gcc make g++ zlib-dev openssl git curl tzdata
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2
RUN golangci-lint --version
WORKDIR /root
RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
WORKDIR /go/src/apple-findmy-to-mqtt
RUN go mod tidy
RUN rm -rf /var/cache/apk/*
RUN chown -R $USER:$USER /go