FROM golang:1.7

ADD . /go/src/github.com/adamwasila/uservices-in-go

RUN go install github.com/adamwasila/uservices-in-go/students

ENTRYPOINT ["/go/bin/students"]

EXPOSE 8080
