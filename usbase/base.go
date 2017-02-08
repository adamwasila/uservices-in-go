package usbase

import (
	"flag"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/nats-io/nats.v1"
	"fmt"
)

func ParseArgs() (string, string, string) {
	var dbUrl *string = flag.String("mongo", "localhost", "Url to Mongo DB")
	var natsUrl *string = flag.String("nats", nats.DefaultURL, "Url to NATS service")
	var restUrl *string = flag.String("bind", "0.0.0.0:8080", "Addres for REST to bind to")

	flag.Parse()
	return *dbUrl, *natsUrl, *restUrl
}

func InitDb(url string) *mgo.Session {
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Println("Can't connect to database.")
		panic(err)
	}
	return session
}

func InitNats(url string) *nats.Conn {
	nc, err := nats.Connect(url)
	if err != nil {
		fmt.Println("Can't connect to NATS service.")
		panic(err)
	}
	return nc
}

func InitRest() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	return router
}
