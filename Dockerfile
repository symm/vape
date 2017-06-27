FROM golang:alpine

RUN mkdir -p /go/src/github.com/symm/vape
COPY . /go/src/github.com/symm/vape

WORKDIR /go/src/github.com/symm/vape

RUN go build
RUN go install

# Final Image
FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=0 /go/bin/vape  .

ENTRYPOINT ["./vape"]
