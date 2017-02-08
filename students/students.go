package main

import (
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/nats-io/nats.v1"
	"github.com/adamwasila/uservices-in-go/usbase"
	"time"
	"encoding/json"
)

type Student struct {
	Id       string `json:"id" binding:"required" bson:"_id"`
	Name     string `json:"name" binding:"required" bson:"firstName"`
	Lastname string `json:"lastname" binding:"required" bson:"lastName"`
	Subjects []string `json:"subjects" binding:"required" bson:"subjects"`
}

type Subject struct {
	Id       string `json:"id" binding:"required" bson:"_id"`
	Name     string `json:"name" binding:"required" bson:"name"`
}

type BaseStudentDto struct {
	Id       string `json:"id" binding:"required" bson:"_id"`
	Name     string `json:"name" binding:"required" bson:"firstName"`
	Lastname string `json:"lastname" binding:"required" bson:"lastName"`
}

type SubjectDto struct {
	Id       string `json:"id" binding:"required" bson:"_id"`
	Name     string `json:"name" binding:"required" bson:"name"`
}

type StudentDto struct {
	*BaseStudentDto
	Subjects []SubjectDto `json:"subjects" binding:"required" bson:"subjects"`
}

var col *mgo.Collection

var natsConnection *nats.Conn

func getStudents(c *gin.Context) {
	var students []Student
	col.Find(bson.M{}).All(&students)

	var studentsDto = []BaseStudentDto{}

	for _, student := range students {
		studentsDto = append(studentsDto, BaseStudentDto {
			Id: student.Id,
			Name: student.Name,
			Lastname: student.Lastname,
		})
	}

	c.JSON(200, studentsDto)
}

func getStudent(c *gin.Context) {
	var id string = c.Param("id")
	var student Student
	var err = col.FindId(id).One(&student)

	if err != nil {
		type MyError struct {
			Id          int `json:"id"`
			Description string `json:"description"`
		}
		var error MyError = MyError{
			0,
			"Error! " + err.Error(),
		}
		c.JSON(400, error)
		return
	}

	var studentDto StudentDto = StudentDto{
		BaseStudentDto: &BaseStudentDto{
			Id: student.Id,
			Name: student.Name,
			Lastname: student.Lastname,
		},
		Subjects: []SubjectDto{},
	}


	if len(student.Subjects) > 0 {
		var subjectsJson, err = json.Marshal(student.Subjects)
		fmt.Printf("nats: %v\n", subjectsJson)
		if err != nil {
			c.String(400, "err: %v", err)
			return
		}

		var msg, err2 = natsConnection.Request("subjects", subjectsJson, 2 * time.Second)
		if err2 != nil {
			c.String(400, "err2: %v", err2)
			return
		}

		var subjects []Subject
		var err3 = json.Unmarshal(msg.Data, &subjects)
		if err3 != nil {
			c.String(400, "err3: %v", err3)
			return
		}

		for _, subject := range subjects {
			studentDto.Subjects = append(studentDto.Subjects, SubjectDto{
				Id: subject.Id,
				Name: subject.Name,
			})
		}

	}

	c.JSON(200, studentDto)
}

func getStudentsNats(msg *nats.Msg) {

}

func initRest(router *gin.Engine)  {
	router.GET("/students", getStudents)
	router.GET("/student/:id", getStudent)
}

func initNats(con *nats.Conn) {
	var _, err = con.Subscribe("subjects", getStudentsNats)
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

	col = session.DB("students").C("students")

	natsConnection = usbase.InitNats(natsUrl)
	initNats(natsConnection)
	defer natsConnection.Close()

	var router = usbase.InitRest()
	initRest(router)
	router.Run(restUrl)
}

