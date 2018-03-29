/*
Package dor is a provider's data collection library and fast HTTP service (build on top of amazing iris framework) which shows a specified domain rank from following providers: Alexa, Majestic, Umbrella OpenDNS and Statvoo

Can be used as a base for a domain categorization / network filters / suspicious domain detection. Data is updated once a day automatically, but it is configurable.


Usage:
	dor-web-inmemory -h
	Usage of dor-web-inmemory:
	-host string
		IP-address to bind (default "127.0.0.1")
	-port string
		Port to bind (default "8080")


Client request example:
	curl 127.0.0.1:8080/rank/github.com

Server response:
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
*/
package dor
