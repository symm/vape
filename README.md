# Vape
[![Build Status](https://img.shields.io/travis/symm/vape.svg)](https://travis-ci.org/symm/vape)
[![Coverage Status](https://img.shields.io/coveralls/symm/vape.svg)](https://coveralls.io/github/symm/vape?branch=master)
[![Docker Pulls](https://img.shields.io/docker/pulls/symm/vape.svg)](https://hub.docker.com/r/symm/vape/)
[![license](https://img.shields.io/github/license/symm/vape.svg)]()

Modern [Smoke testing](https://en.wikipedia.org/wiki/Smoke_testing) tool written in Go. Inspired by [Shisha](https://github.com/namshi/shisha)

![Success](/assets/success.png?raw=true "Success")
![Failure](/assets/failure.png?raw=true "Failure")

# How to use

## As a binary

Grab a binary from our [Releases page](https://github.com/symm/vape/releases) or build one by checking out this repo and running `make`

Then create a `Vapefile` file in the format:
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

then execute `vape http://your.domain` to run the tests

## As a container

No need to download binaries or compile the project, we publish a ready made image on [Docker Hub](https://hub.docker.com/r/symm/vape/)

Just create the `Vapefile` file as above and mount it inside a container:

```shell
docker run \
    --rm \
    -t \
    -v $(PWD)/Vapefile.example:/Vapefile \
    symm/vape:latest \
    https://your.domain
```

### Optional Arguments

`-config full/path/to/Vapefile`: specify an alternative to looking for `Vapefile` in the current directory
`-skip-ssl-verification`: Ignore bad / self signed SSL certificates

## TODO

This project is HackDay™ quality. In need of test coverage and refactoring
