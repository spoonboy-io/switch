# Switch

### Preprocess JSON to `name` and `value` keys for Morpheus 

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/spoonboy-io/switch?style=flat-square)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/spoonboy-io/switch?style=flat-square)](https://goreportcard.com/report/github.com/spoonboy-io/switch)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/spoonboy-io/switch/build.yml?branch=master&style=flat-square)](https://github.com/spoonboy-io/switch/actions/workflows/build.yml)
[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/spoonboy-io/switch/unit_test.yml?branch=master&label=tests&style=flat-square)](https://github.com/spoonboy-io/switch/actions/workflows/unit_test.yml)
[ ![ GitHub Release Date ](https://img.shields.io/github/release-date/spoonboy-io/switch?style=flat-square) ](https://github.com/spoonboy-io/switch/releases)
[![GitHub commits since latest release (by date)](https://img.shields.io/github/commits-since/spoonboy-io/switch/latest?style=flat-square)](https://github.com/spoonboy-io/switch/commits)
[![GitHub](https://img.shields.io/github/license/spoonboy-io/switch?label=license&style=flat-square)](LICENSE)


## About

Switch performs preprocessing on JSON data. Designed to simplify Morpheus option list creation and management,
Switch can fetch a target URL and parse it, creating 
`name` and `value` keys in a new JSON file according to the YAML configuration rules provided.

Switch can manage multiple source files, and, using a Time to Live (TTL) for each source it will refresh the data
at the specified interval.

Switch saves the data as a new JSON file in a location you specify. Any web server can serve it, you could also use
the [Dujour JSON/CSV data file server](https://github.com/spoonboy-io/dujour).


## Why?

Formulating the JavaScript for Morpheus Option List REST translation scripts can be complex, while any JSON file which presents
as an array of `name` and `value` keys, needs no translation script whatsoever since Morpheus is able to interpret it automatically.

So, Switch takes a complex JSON payload, which would require a translation script to parse in Morpheus, and creates a
simple representation of the data needed for the option list.

Basic use, requires no translation script, more complicated use (such as dependent inputs) will mean any translation 
script is much simplified.

## Releases

You can find the [latest software here](https://github.com/spoonboy-io/switch/releases/latest).

## Usage

Switch will look for and parse a [sources.yaml file](https://github.com/spoonboy-io/switch/blob/master/sources.yaml) which has the following format:

```yaml
---
- source:
    description: Some blog posts (array)
    url: https://jsonplaceholder.typicode.com/posts
    token:
    extract:
      root:
      name: title
      value: id
    ttl: 5 # minutes
    save:
      folder: test
      filename: test.json
```

The above configuration creates this file `./test/test.json` locally:

```json

[{
	"name": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
	"value": 1
}, {
	"name": "qui est esse",
	"value": 2
}, {
	"name": "ea molestias quasi exercitationem repellat qui ipsa sit aut",
	"value": 3
}, {
	"name": "eum et est occaecati",
	"value": 4
}, {
	"name": "nesciunt quas odio",
	"value": 5
},
  
  // ... abbreviated
  
  {
	"name": "at nam consequatur ea labore ea harum",
	"value": 100
}]
```

A similar config, which adds a root of "batter", can parse remote JSON like this:

```json
{
  "items": {
    "item": [
      {
        "batters": {
          "batter": [
            {
              "id": "1001",
              "type": "Regular"
            },
            {
              "id": "1002",
              "type": "Chocolate"
            },
            {
              "id": "1003",
              "type": "Blueberry"
            },
            {
              "id": "1004",
              "type": "Devil's Food"
            }
          ]
        },
        "id": "0001",
        "name": "Cake",
        "ppu": 0.55,
        "topping": [
          {
            "id": "5001",
            "type": "None"
          },
          {
            "id": "5002",
            "type": "Glazed"
          },
          {
            "id": "5005",
            "type": "Sugar"
          },
          {
            "id": "5007",
            "type": "Powdered Sugar"
          },
          {
            "id": "5006",
            "type": "Chocolate with Sprinkles"
          },
          {
            "id": "5003",
            "type": "Chocolate"
          },
          {
            "id": "5004",
            "type": "Maple"
          }
        ],
        "type": "donut"
      }
    ]
  }
}
```

To this:

```json
[{
	"name": "Regular",
	"value": "1001"
}, {
	"name": "Chocolate",
	"value": "1002"
}, {
	"name": "Blueberry",
	"value": "1003"
}, {
	"name": "Devil's Food",
	"value": "1004"
}]
```

### Installation
Grab the tar.gz or zip archive for your OS from the [releases page](https://github.com/spoonboy-io/switch/releases/latest).

Unpack it to the target host, and then start the application.

```
./switch
```

Or with nohup..

```
nohup ./switch &
```

### TODO

- Unit tests only cover the extraction routines
- Only manually tested with unauthenticated GET requests at this time
- The extraction code, caters for simple arrays and objects keys which store an array

### License
Licensed under [Mozilla Public License 2.0](LICENSE)
