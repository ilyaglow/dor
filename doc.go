/*
Package dor is a domain rank data collection library and fast HTTP service
which shows a specified domain's rank from the following providers:
* Alexa
* Majestic
* Umbrella OpenDNS
* Statvoo
* Open PageRank
* Tranco
* Quantcast

Can be used as a base for a domain categorization, network filters or
suspicious domain detection. Data is updated automatically by dor-insert once a
day by default.

See service/dor-http/dor.go for an example of the Dor HTTP service and
cmd/dor-insert/dor-insert.go for the data insertion script.

Client request example:
	curl 127.0.0.1:8080/rank/github.com

Server response:
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
		},
		{
		  "domain": "github.com",
		  "rank": 114,
		  "date": "2019-05-04T00:00:00Z",
		  "source": "quantcast",
		  "raw": ""
		}
	  ],
	  "timestamp": "2018-01-11T18:07:09.186271429Z"
	}
*/
package dor
