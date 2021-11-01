[![run-tests](https://github.com/DionTech/services-lfu-cache/actions/workflows/go.yml/badge.svg)](https://github.com/DionTech/services-lfu-cache/actions/workflows/go.yml)
[![Twitter Follow](https://img.shields.io/twitter/follow/dion_tech?style=social)](https://twitter.com/dion_tech)
# About
## Project Description
This is a project to learn some elements of go. It is a thread-safe cache algorithm, based on the LFU principle. Also you can define the max heap size the runtime of this cache service may use. The aim is to not include any other package dependencies away from the standard go library to can really learn all the basics. 

## Architecture
At the moment, nothing to see here :) Later it will use docker to can be integrated as a part-based tcp service to get / set cache elements.

# Setup
## Install

```sh
go mod tidy
```

## Test 
```sh
go test ./... -cover 
```

## Benchmarks


