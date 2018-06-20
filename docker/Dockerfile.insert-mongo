FROM golang:alpine AS build-env
LABEL maintainer "Ilya Glotov <ilya@ilyaglotov.com>" \
      repository "https://github.com/ilyaglow/dor"

ENV CGO_ENABLED 0

RUN apk -U --no-cache add git \
  && go get -v -ldflags="-s -w" github.com/ilyaglow/dor/cmd/dor-insert-mongo \
  && apk del git

FROM alpine:edge

RUN apk -U --no-cache add ca-certificates \
  && adduser -D app

COPY --from=build-env /go/bin/dor-insert-mongo /dor-insert

USER app

ENTRYPOINT ["/dor-insert"]
