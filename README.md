# Vape

Modern [Smoke testing](https://en.wikipedia.org/wiki/Smoke_testing) tool written in Go. Inspired by [Shisha](https://github.com/namshi/shisha)

# How to use

## As a binary

Create a `Vapefile` file in the format:
```json
[
  {
    "uri": "/health",
    "expectedStatusCode": 200
  },
  {
    "uri": "/page-that-should-not-exist",
    "expectedStatusCode": 404
  }
]
```

then execute `vape http://your.domain` to run the checks

## As a container

Create the `Vapefile` file as above but be sure to mount it inside the container:

```bash
docker run \
    --rm \
    -t \
    -v $(PWD)/Vapefile.example:/Vapefile \
    symm/vape:latest
```

## TODO

This project is HackDay™ quality. In need of test coverage and refactoring
