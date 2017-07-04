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

Create a `Vapefile` file in the format:
```json
[
  {
    "uri": "/status/418",
    "expected_status_code": 418,
    "content": "teapot"
  },
  {
    "uri": "/status/200",
    "expected_status_code": 200
  },
  {
    "uri": "/status/304",
    "expected_status_code": 304
  },
  {
    "uri": "/status/500",
    "expected_status_code": 500
  }
]
```

then execute `vape http://your.domain` to run the tests

## As a container

Create the `Vapefile` file as above but be sure to mount it inside the container:

```shell
docker run \
    --rm \
    -t \
    -v $(PWD)/Vapefile.example:/Vapefile \
    symm/vape:latest \
    https://your.domain
```

### Optional Arguments

`-config full/path/to/Vapefile` specify an alternative to looking for `Vapefile` in the current directory

## TODO

This project is HackDay™ quality. In need of test coverage and refactoring
