# Microservices in go

Simple, fun-type, project started to experiment with (and learn) golang and to test how to write small&quick microservices efficiently in this setup.

No practical use. No ambitions to be beautiful code example either. 

Quick start:

* install go, glide, docker, docker-compose; only first one is obligatory
* run glide install to install all dependencies to vendor directory
* (A) for in-container build, run:
```
docker-compose up -d
```
* (B) for regular build, run:
```
go install github.com/adamwasila/uservices-in-go/students
go install github.com/adamwasila/uservices-in-go/subjects
```

Libraries used:

* gin - http framework/router to provide REST API
* mgo - MongoDB driver to get data from a DB
* nats - client for nats messaging to connect all microservice instances together

Other:

* MongoDB
* Docker (optional)
