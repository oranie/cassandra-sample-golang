package main

import (
	"./internal/pkg"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"log"
	"net/http"
	"time"
)

type Chat struct {
	Name     string
	Time     int64
	Chatroom string
	Comment  string
}

func main() {
	env, session, chatData := initApp()
	defer session.Close()

	log.Printf("session : %v", session)

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
		if &chatData != nil {
			chatData = chat.GenerateChatData()
		}
		//insert test data
		chat.InsertData(session, &chatData)
		//select insert data
		chat.SelectTestData(session, &chatData)
		//select all data at chat table
		result := chat.AllSelectData(session)
		c.String(http.StatusOK, "Test done.", result)
	})

	r.StaticFile("/chat", "./web/livechat.html")

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

	r.GET("/chat/comments/latest", func(c *gin.Context) {
		chatroom := "game_room-oranie"
		chatData := chat.ChatroomLatestData(session, chatroom)
		json, err := json.Marshal(chatData)
		if err != nil {
			panic(err)
		}
		c.String(http.StatusOK, string(json))
	})

	r.GET("/chat/comments/all", func(c *gin.Context) {
		chatroom := "game_room-oranie"
		chatData := chat.ChatroomAllData(session, chatroom)
		json, err := json.Marshal(chatData)
		if err != nil {
			panic(err)
		}
		c.String(http.StatusOK, string(json))
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
	portString := ":" + env.AppPort
	r.Run(portString)
}

func initApp() (chat.Env, *gocql.Session, chat.Chat) {
	env := chat.GetEnvValue()

	if env.AppEnv == "prd" || env.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	session, error := chat.CreateCassandraSession()
	if error != nil {
		fmt.Println(error)
	}

	//check and create chat table
	chat.CreateChatTable(session)

	//generate test data
	chatData := chat.GenerateChatData()
	return env, session, chatData
}
