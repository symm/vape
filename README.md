# Vape

Modern [Smoke testing](https://en.wikipedia.org/wiki/Smoke_testing) tool written in Go. Inspired by [Shisha](https://github.com/namshi/shisha)

# How to use

## As a binary

Create a `.smoke` file in the format:
```
https://httpbin.org/status/418 418
https://httpbin.org/status/200 200
https://httpbin.org/status/500 500
https://httpbin.org/status/404 404
```

then execute the `vape` binary to run the checks

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
