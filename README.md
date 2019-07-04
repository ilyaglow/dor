[![Build Status](https://travis-ci.org/ilyaglow/dor.svg?branch=master)](https://travis-ci.org/ilyaglow/dor)
[![](https://godoc.org/github.com/ilyaglow/dor?status.svg)](http://godoc.org/github.com/ilyaglow/dor)

DOR - Domain Ranker
-------------------

Fast HTTP service which shows a specified domain rank from following providers:
- [Alexa](https://www.alexa.com/topsites)
- [Majestic](https://blog.majestic.com/development/alexa-top-1-million-sites-retired-heres-majestic-million/)
- [Umbrella OpenDNS](https://umbrella.cisco.com/blog/2016/12/14/cisco-umbrella-1-million/)
- [Open PageRank](https://www.domcop.com/top-10-million-domains)
- [Tranco List](https://tranco-list.eu/)
- [Quantcast](https://www.quantcast.com/top-sites/)

Can be used as a base for a domain categorization / network filters /
suspicious domain detection.

Data is updated once a day automatically.

Supported types of storages:
* Clickhouse (recommended)
* MongoDB
* In-Memory

You can easily add the storage you like by implementing _Storage_ interface.

## Installation

Check out the [releases page](https://github.com/ilyaglow/dor/releases).

## HTTP service usage

Use Clickhouse storage located at `clickhouse` and bind to port `8080`:
```
go run service/dor-http/dor.go \
    -storage=clickhouse \
    -storage-url=tcp://clickhouse:9000 \
    -listen-addr=:8080
```

## Fill database with the data

```
go run cmd/dor-insert/dor-insert \
    -storage=clickhouse \
    -storage-url=tcp://clickhouse:9000
```

## Docker usage

Project has [docker-compose](docker-compose.yml) that uses Clickhouse as a
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
      "rank": 40,
      "date": "2019-06-01T00:00:00Z",
      "source": "tranco",
      "raw": ""
    },
    {
      "domain": "github.com",
      "rank": 123,
      "date": "2019-06-01T00:00:00Z",
      "source": "quantcast",
      "raw": ""
    },
    {
      "domain": "github.com",
      "rank": 2555,
      "date": "2019-06-01T00:00:00Z",
      "source": "umbrella",
      "raw": ""
    },
    {
      "domain": "github.com",
      "rank": 45,
      "date": "2019-06-01T00:00:00Z",
      "source": "alexa",
      "raw": ""
    },
    {
      "domain": "github.com",
      "rank": 29,
      "date": "2019-06-01T00:00:00Z",
      "source": "majestic",
      "raw": "29,23,github.com,com,179597,516432,github.com,com,29,23,179410,516889"
    },
    {
      "domain": "github.com",
      "rank": 39,
      "date": "2019-06-01T00:00:00Z",
      "source": "pagerank",
      "raw": ""
    }
  ],
  "timestamp": "2019-06-01T22:06:12.662492307Z"
}
```
