FROM golang:1.9-alpine3.6

ADD . /go/src/github.com/oemdaro/mqtt-microservices-example/auth-service

RUN go install github.com/oemdaro/mqtt-microservices-example/auth-service

ENTRYPOINT ["/go/bin/auth-service"]

EXPOSE 50051
