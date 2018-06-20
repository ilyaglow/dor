FROM golang:alpine AS build-env
LABEL maintainer "Ilya Glotov <ilya@ilyaglotov.com>" \
      repository "https://github.com/ilyaglow/dor"

ENV CGO_ENABLED 0

RUN apk -U --no-cache add git \
  && go get -v -ldflags="-s -w" github.com/ilyaglow/dor/service/dor-web-mongodb \
  && apk del git

FROM alpine:edge

RUN apk -U --no-cache add ca-certificates \
  && adduser -D app

COPY --from=build-env /go/bin/dor-web-mongodb /dor

USER app

ENTRYPOINT ["/dor"]
