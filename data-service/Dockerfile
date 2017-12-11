FROM golang:1.9-alpine3.6

ADD . /go/src/github.com/oemdaro/mqtt-microservices-example/data-service

RUN go install github.com/oemdaro/mqtt-microservices-example/data-service

ENTRYPOINT ["/go/bin/data-service"]
