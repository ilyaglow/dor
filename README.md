[![Build Status](https://travis-ci.org/ilyaglow/dor.svg?branch=master)](https://travis-ci.org/ilyaglow/dor)
[![](https://godoc.org/github.com/ilyaglow/dor?status.svg)](http://godoc.org/github.com/ilyaglow/dor)

DOR - Domain Ranker
-------------------

Fast HTTP service which shows a specified domain rank from following providers:
- [Alexa](https://www.alexa.com/topsites)
- [Majestic](https://blog.majestic.com/development/alexa-top-1-million-sites-retired-heres-majestic-million/)
- [Umbrella OpenDNS](https://umbrella.cisco.com/blog/2016/12/14/cisco-umbrella-1-million/)
- [Statvoo](https://statvoo.com/top/sites)
- [Open PageRank](https://www.domcop.com/top-10-million-domains)

Can be used as a base for a domain categorization / network filters /
suspicious domain detection.

Data is updated once a day automatically.

Right now only in-memory and MongoDB storages are supported, but _Dor_ was
built with flexibility in mind, so you can easily add the storage you like by
implementing _Storage_ interface.

## Installation

Check out the [releases page](https://github.com/ilyaglow/dor/releases).

### Manual build

**dor** supports **Go 1.9 and later**

```
go get -u github.com/ilyaglow/dor
go install ./...
```

## Web service usage

Use MongoDB storage located at `mongoserver` and bind to port `8080`
```
DOR_MONGO_URL=mongoserver DOR_PORT=8080 dor-web-mongodb
```

## Fill database with the data

```
DOR_MONGO_URL=mongoserver go run cmd/dor-insert-mongo/dor-insert-mongo
```

Or if you want just in-memory database:
```
dor-web-inmemory -h

Usage of dor-web-inmemory:
  -listen string
    	Listen address to bind (default "127.0.0.1:8080")
```

## Docker usage

Project has [docker-compose](docker-compose.yml) that uses MongoDB as a
storage. Make changes here accordingly if any (folder for data persistence,
ports etc).

```
docker-compose up -d
```


## Client usage

```sh
$: curl 127.0.0.1:8080/rank/github.com

{
  "data": "github.com",
  "ranks": [
    {
      "domain": "github.com",
      "rank": 33,
      "last_update": "2018-01-11T18:01:27.251103268Z",
      "source": "majestic"
    },
    {
      "domain": "github.com",
      "rank": 66,
      "last_update": "2018-01-11T18:01:27.97067767Z",
      "source": "statvoo"
    },
    {
      "domain": "github.com",
      "rank": 72,
      "last_update": "2018-01-11T18:04:26.267833256Z",
      "source": "alexa"
    },
    {
      "domain": "github.com",
      "rank": 2367,
      "last_update": "2018-01-11T18:06:50.866600102Z",
      "source": "umbrella"
    },
    {
      "domain": "github.com",
      "rank": 115,
      "last_update": "2018-03-27T17:01:13.535Z",
      "source": "pagerank"
    }
  ],
  "timestamp": "2018-01-11T18:07:09.186271429Z"
}
```
