FROM golang:alpine

RUN apk update && apk add git
RUN go get -u github.com/golang/dep/cmd/dep

RUN mkdir -p /go/src/github.com/symm/vape
COPY . /go/src/github.com/symm/vape

WORKDIR /go/src/github.com/symm/vape

RUN dep ensure
RUN go build
RUN go install

# Final Image
FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=0 /go/bin/vape  .

CMD ["./vape"]
