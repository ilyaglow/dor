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

## Usage

```
dor-webservice -h

Usage of dor-webservice:
  -host string
    	IP-address to bind (default "127.0.0.1")
  -port string
    	Port to bind (default "8080")
```
