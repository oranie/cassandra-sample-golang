package main

import (
	"./service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Chat struct {
	Name     string
	Time     int64
	Chatroom string
	Comment  string
}

func main() {
	session, error := service.CreateCassandraSession()
	if error != nil {
		fmt.Println(error)
	}

	defer session.Close()
	service.CreateChatTable(session)

	//generate test data
	chatData := service.GenerateChatData()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "server status id good",
		})
	})

	r.GET("/run-test", func(c *gin.Context) {
		//insert test data
		service.InsertData(session, &chatData)
		//select insert data
		service.SelectTestData(session, &chatData)
		//select all data at chat table
		service.AllSelectData(session)
		c.String(http.StatusOK, "Test done.")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/insertstatus", func(c *gin.Context) {
		chatData := service.SelectTestData(session, &chatData)
		json, err := json.Marshal(chatData)
		if err != nil {
			panic(err)
		}
		c.String(http.StatusOK, string(json))
	})
	r.GET("/alldata", func(c *gin.Context) {
		allChatData := service.AllSelectData(session)
		json, err := json.Marshal(allChatData)
		if err != nil {
			panic(err)
		}
		c.String(http.StatusOK, string(json))
	})

	r.Run()
}
