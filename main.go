package main

import (
	"./internal/pkg"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

type Chat struct {
	Name     string
	Time     int64
	Chatroom string
	Comment  string
}

func main() {
	env := os.Getenv("ENV")
	if env == "prd" || env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	session, error := chat.CreateCassandraSession()
	if error != nil {
		fmt.Println(error)
	}

	defer session.Close()
	chat.CreateChatTable(session)

	//generate test data
	chatData := chat.GenerateChatData()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "server status id good",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/run-test", func(c *gin.Context) {
		//insert test data
		chat.InsertData(session, &chatData)
		//select insert data
		chat.SelectTestData(session, &chatData)
		//select all data at chat table
		chat.AllSelectData(session)
		c.String(http.StatusOK, "Test done.")
	})

	r.POST("/chat/comments/add", func(c *gin.Context) {
		var json chat.Chat
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		postData := chat.Chat{
			Name:     json.Name,
			Time:     time.Now().UnixNano(),
			Chatroom: "oranie-room",
			Comment:  json.Comment,
		}

		fmt.Printf("%s", json)
		resp := chat.InsertData(session, &postData)

		c.String(http.StatusOK, resp)
	})

	r.GET("/insertstatus", func(c *gin.Context) {
		chatData := chat.SelectTestData(session, &chatData)
		json, err := json.Marshal(chatData)
		if err != nil {
			panic(err)
		}
		c.String(http.StatusOK, string(json))
	})

	r.GET("/alldata", func(c *gin.Context) {
		allChatData := chat.AllSelectData(session)
		json, err := json.Marshal(allChatData)
		if err != nil {
			panic(err)
		}
		c.String(http.StatusOK, string(json))
	})

	r.Run()
}
