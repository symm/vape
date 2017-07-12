FROM golang:alpine

LABEL maintainer="me@gazj.co.uk"

RUN mkdir -p /go/src/github.com/symm/vape
COPY . /go/src/github.com/symm/vape

WORKDIR /go/src/github.com/symm/vape

RUN go build
RUN go install

# Final Image
FROM alpine:3.6
RUN apk add --no-cache ca-certificates
COPY --from=0 /go/bin/vape  .

ENTRYPOINT ["./vape"]
