[![Build Status](https://travis-ci.org/ilyaglow/dor.svg?branch=master)](https://travis-ci.org/ilyaglow/dor) [![](https://godoc.org/github.com/ilyaglow/dor?status.svg)](http://godoc.org/github.com/ilyaglow/dor)

DOR - Domain Ranker
-------------------

Fast HTTP service (build on top of amazing [iris framework](https://github.com/kataras/iris)) which shows a specified domain rank from following providers:
- [Alexa](https://www.alexa.com/topsites)
- [Majestic](https://blog.majestic.com/development/alexa-top-1-million-sites-retired-heres-majestic-million/)
- [Umbrella OpenDNS](https://umbrella.cisco.com/blog/2016/12/14/cisco-umbrella-1-million/)
- [Statvoo](https://statvoo.com/top/sites)

Can be used as a base for a domain categorization / network filters / suspicious domain detection.

Data is updated once a day automatically.

## Installation

Check out the [releases page](https://github.com/ilyaglow/dor/releases).

### Manual build

**dor** supports **Go 1.9 and later**

```
go get -u github.com/ilyaglow/dor/cmd/dor-webservice
```

## Web service usage

```
dor-webservice -h

Usage of dor-webservice:
  -host string
    	IP-address to bind (default "127.0.0.1")
  -port string
    	Port to bind (default "8080")
```

## Client usage

```
$: curl 127.0.0.1:8080/rank/github.com

{
  "data": "github.com",
  "ranks": [
    {
      "domain": "github.com",
      "rank": 33,
      "last_update": "2018-01-11T18:01:27.251103268Z",
      "description": "majestic"
    },
    {
      "domain": "github.com",
      "rank": 66,
      "last_update": "2018-01-11T18:01:27.97067767Z",
      "description": "statvoo"
    },
    {
      "domain": "github.com",
      "rank": 72,
      "last_update": "2018-01-11T18:04:26.267833256Z",
      "description": "alexa"
    },
    {
      "domain": "github.com",
      "rank": 2367,
      "last_update": "2018-01-11T18:06:50.866600102Z",
      "description": "umbrella"
    }
  ],
  "timestamp": "2018-01-11T18:07:09.186271429Z"
}
```
