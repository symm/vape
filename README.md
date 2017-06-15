# Vape

Modern [Smoke testing](https://en.wikipedia.org/wiki/Smoke_testing) tool written in Go. Inspired by [Shisha](https://github.com/namshi/shisha)

# How to use

## As a binary

Create a `Vapefile` file in the format:
```
[
  {
    "uri": "/status/200",
    "expectedStatusCode": 200
  },
  {
    "uri": "/status/500",
    "expectedStatusCode": 500
  }
]
```

then execute `vape http://your.domain/` to run the checks

## As a container

Create the smoke file as above but be sure to mount it inside the container:

```
docker run \
    --rm \
    -t \
    -v $(PWD)/smoke.example:/.smoke \
    symm/vape:latest
```

## TODO

This project is HackDay™ quality. In need of test coverage and refactoring
