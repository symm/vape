# Vape
[![Build Status](https://img.shields.io/travis/symm/vape.svg)](https://travis-ci.org/symm/vape)
[![Codecov](https://img.shields.io/codecov/c/github/symm/vape.svg)](https://codecov.io/gh/symm/vape)
[![Docker Pulls](https://img.shields.io/docker/pulls/symm/vape.svg)](https://hub.docker.com/r/symm/vape/)
[![license](https://img.shields.io/github/license/symm/vape.svg)]()

Modern [Smoke testing](https://en.wikipedia.org/wiki/Smoke_testing) tool written in Go. Inspired by [Shisha](https://github.com/namshi/shisha)

Vape is intended to be used within a [Continuous Delivery pipeline](https://en.wikipedia.org/wiki/Continuous_delivery)
as a post-deployment step to quickly verify if a release is good or not.

It can quickly make assertions about the status code and content for a list of URIs on a given server.

![Success](/assets/success.png?raw=true "Success")
![Failure](/assets/failure.png?raw=true "Failure")

# How to use

## Create a config file

Create a file named `Vapefile` file in the format:
```json
[
  {
    "uri": "/status/418",
    "status_code": 418,
    "content": "teapot"
  },
  {
    "uri": "/status/200",
    "status_code": 200
  },
  {
    "uri": "/status/304",
    "status_code": 304
  },
  {
    "uri": "/status/500",
    "status_code": 500
  }
]
```

The `uri` and `status_code` are required, `content` check is optional

## Run the app from a container (Recommended)

We publish a ready made image on [Docker Hub](https://hub.docker.com/r/symm/vape/)

Just create the `Vapefile` file as above and mount it inside a container:

```bash
docker run \
    --rm \
    -t \
    -v $(PWD)/Vapefile.example:/Vapefile \
    symm/vape:latest \
    https://your.domain
```

## Run the app from a binary

Grab a binary from our [releases page](https://github.com/symm/vape/releases) or build one by checking out this repo and running `make`
then execute `./vape http://your.domain` to run the tests


## Optional flags

The following optional command line flags may be passed:

```bash
Usage of ./vape:
  -concurrency int
    	The maximum number of requests to make at a time (default 3)
  -config string
    	The full path to the Vape configuration file (default "Vapefile")
  -skip-ssl-verification
    	Ignore bad SSL certs
```

For example:

```bash
./vape -concurrency 10 -config vape.conf -skip-ssl-verification http://httpbin.org
```
