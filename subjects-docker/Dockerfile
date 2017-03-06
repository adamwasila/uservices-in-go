FROM golang:1.7

ADD . /go/src/github.com/adamwasila/uservices-in-go

RUN go install github.com/adamwasila/uservices-in-go/subjects

ENTRYPOINT ["/go/bin/subjects"]

EXPOSE 8080
