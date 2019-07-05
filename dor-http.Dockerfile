FROM golang:alpine AS build-env
LABEL maintainer "Ilya Glotov <ilya@ilyaglotov.com>" \
      repository "https://github.com/ilyaglow/dor"

ENV CGO_ENABLED=0 \
    GO111MODULE=on

COPY . /go/src/dor

RUN apk -U --no-cache add git \
  && cd /go/src/dor \
  && go build -ldflags="-s -w" -o /dor service/dor-http/dor.go \
  && apk del git

FROM alpine:edge

RUN apk -U --no-cache add ca-certificates \
  && adduser -D app

COPY --from=build-env /dor /dor

USER app

ENTRYPOINT ["/dor"]
