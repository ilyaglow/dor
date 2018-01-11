[![](https://godoc.org/github.com/ilyaglow/dor?status.svg)](http://godoc.org/github.com/ilyaglow/dor)

DOR - Domain Ranker
-------------------

Fast HTTP service which shows a specified domain rank from following providers:
- Alexa
- Majestic
- Umbrella OpenDNS
- Statvoo

Data is updated once a day automatically

## Installation

```
go get -u github.com/ilyaglow/dor
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
