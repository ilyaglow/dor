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
- [Tranco List](https://tranco-list.eu/)

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

## Web service usage

Use Clickhouse storage located at `clickhouse` and bind to port `8080`
```
DOR_STORAGE=clickhouse DOR_STORAGE_URL=tcp://clickhouse:9000 DOR_PORT=8080 dor-web
```

## Fill database with the data

```
DOR_STORAGE_URL=tcp://clickhouse:9000 DOR_STORAGE=clickhouse go run cmd/dor-insert/dor-insert
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
      "rank": 33,
      "date": "2018-01-11T18:01:27.251103268Z",
      "source": "majestic",
      "raw": "29,23,github.com,com,179825,518189,github.com,com,29,23,179994,518726"
    },
    {
      "domain": "github.com",
      "rank": 66,
      "date": "2018-01-11T18:01:27.97067767Z",
      "source": "statvoo",
      "raw": ""
    },
    {
      "domain": "github.com",
      "rank": 72,
      "date": "2018-01-11T18:04:26.267833256Z",
      "source": "alexa",
      "raw": ""
    },
    {
      "domain": "github.com",
      "rank": 2367,
      "last_update": "2018-01-11T18:06:50.866600102Z",
      "source": "umbrella",
      "raw": ""
    },
    {
      "domain": "github.com",
      "rank": 115,
      "last_update": "2018-03-27T17:01:13.535Z",
      "source": "pagerank",
      "raw": ""
    },
    {
      "domain": "github.com",
      "rank": 68,
      "last_update": "2018-03-27T17:01:13.535Z",
      "source": "tranco",
      "raw": ""
    }
  ],
  "timestamp": "2018-01-11T18:07:09.186271429Z"
}
```
