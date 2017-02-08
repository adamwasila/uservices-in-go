package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/nats-io/nats.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/adamwasila/uservices-in-go/usbase"
	"encoding/json"
	"fmt"
)

type Subject struct {
	Id       string `json:"id" binding:"required" bson:"_id"`
	Name     string `json:"name" binding:"required" bson:"name"`
}

var col *mgo.Collection

var natsConnection *nats.Conn

func getSubjects(c *gin.Context) {
	var students []Subject
	col.Find(bson.M{}).All(&students)
	c.JSON(200, students)
}

func getSubject(c *gin.Context) {
	var id = c.Param("id")
	var student Subject
	var err = col.FindId(id).One(&student)

	if err == nil {
		c.JSON(200, student)
	} else {
		type MyError struct {
			Id          int `json:"id"`
			Description string `json:"description"`
		}
		var error MyError = MyError{
			0,
			"Error!",
		}
		c.JSON(400, error)
	}
}

func getSubjectsNats(msg *nats.Msg) {
	var subjectIds []string
	var subjects []Subject

	var err = json.Unmarshal(msg.Data, &subjectIds)
	if err != nil {
		fmt.Println("Failing to parse nats request...")
		var data, _ = json.Marshal([]string{})
		natsConnection.Publish(msg.Reply, data)
		return
	}

	for _, subjectId := range subjectIds {
		var subject Subject
		col.FindId(subjectId).One(&subject)
		subjects = append(subjects, subject)
	}

	var data, _ = json.Marshal(subjects)
	natsConnection.Publish(msg.Reply, data)
}

func initRest(router *gin.Engine)  {
	router.GET("/subjects", getSubjects)
	router.GET("/subject/:id", getSubject)
}

func initNats(con *nats.Conn) {
	var _, err = con.Subscribe("subjects", getSubjectsNats)
	if err != nil {
		fmt.Println("Error while subscribing to...")
		panic(err)
	}

}

func main() {
	fmt.Println("Start...")
	var dbUrl, natsUrl, restUrl = usbase.ParseArgs()

	var session = usbase.InitDb(dbUrl)
	defer session.Close()

	col = session.DB("subjects").C("subjects")

	natsConnection = usbase.InitNats(natsUrl)
	initNats(natsConnection)
	defer natsConnection.Close()

	var router = usbase.InitRest()
	initRest(router)
	router.Run(restUrl)
}
