# Vape [![Build Status](https://travis-ci.org/symm/vape.svg?branch=master)](https://travis-ci.org/symm/vape) [![Coverage Status](https://coveralls.io/repos/github/symm/vape/badge.svg?branch=master)](https://coveralls.io/github/symm/vape?branch=master)

Modern [Smoke testing](https://en.wikipedia.org/wiki/Smoke_testing) tool written in Go. Inspired by [Shisha](https://github.com/namshi/shisha)

# How to use

## As a binary

Create a `Vapefile` file in the format:
```json
[
  {
    "uri": "/health",
    "expected_status_code": 200
  },
  {
    "uri": "/page-that-should-not-exist",
    "expected_status_code": 404
  }
]
```

then execute `vape http://your.domain` to run the checks

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
