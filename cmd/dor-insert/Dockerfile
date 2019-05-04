FROM golang:alpine AS build-env
LABEL maintainer "Ilya Glotov <ilya@ilyaglotov.com>" \
      repository "https://github.com/ilyaglow/dor"

ENV CGO_ENABLED 0

COPY dor-insert.go /go/src/dor-insert/dor-insert.go

RUN apk -U --no-cache add git \
  && cd /go/src/dor-insert \
  && go get -v . \
  && go build -ldflags="-s -w" -o /dor-insert \
  && apk del git

FROM alpine:edge

RUN apk -U --no-cache add ca-certificates \
  && adduser -D app

COPY --from=build-env /dor-insert /dor-insert

USER app

ENTRYPOINT ["/dor-insert"]
