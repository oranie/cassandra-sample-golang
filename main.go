package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/oranie/cassandra-sample-golang/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Comments struct {
	Response []chat.Comment `form:"name" json:"response"`
}

func main() {
	env, session, chatData := chat.InitApp()
	defer session.Close()

	log.Printf("session : %v", session)

	r := gin.Default()
	ApiEndpoint := "http://localhost:" + env.AppPort + "/"
	// local env check
	localCheck := regexp.MustCompile(`localhost|127.0.0.1`)
	if !localCheck.MatchString(env.AppEndpoint) {
		ApiEndpoint = "https://" + env.AppEndpoint + "/"
	}
	log.Println("App Endpoint : ", ApiEndpoint)

	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "server status id good",
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

	r.GET("/chat", func(c *gin.Context) {

		f, err := os.Open("./web/livechat.html")
		if err != nil {
			fmt.Println("file read error", err)
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal("html file read error: ", err)
		}

		//html static http://localhost:8080/
		// if public endpoint
		ApiEndpoint := "http://localhost:" + env.AppPort + "/"
		rep := regexp.MustCompile(`http://localhost:8080/`)
		// local env check
		localCheck := regexp.MustCompile(`localhost|127.0.0.1`)
		if !localCheck.MatchString(env.AppEndpoint) {
			ApiEndpoint = "https://" + env.AppEndpoint + "/"
		}

		str := rep.ReplaceAllString(string(b), ApiEndpoint)
		c.Header("Content-Type", "text/html")
		c.Header("Access-Control-Allow-Origin", "*")

		c.String(http.StatusOK, str)
	})

	r.POST("/chat/comments/add", func(c *gin.Context) {
		var json chat.Comment
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		postData := chat.Comment{
			Name:     json.Name,
			Time:     time.Now().UnixNano(),
			Chatroom: "game_room-oranie",
			Comment:  json.Comment,
		}

		fmt.Printf("%v", json)
		resp := chat.InsertData(session, &postData)

		c.JSON(http.StatusOK, resp)
	})

	r.GET("/chat/comments/latest", func(c *gin.Context) {
		chatroom := "game_room-oranie"
		chatData := chat.ChatroomLatestData(session, chatroom)
		comnents := Comments{Response: chatData}

		log.Println("latest data :", comnents)
		c.JSON(http.StatusOK, comnents)
	})

	r.GET("/chat/comments/all", func(c *gin.Context) {
		chatroom := "game_room-oranie"
		chatData := chat.ChatroomAllData(session, chatroom)
		comnents := Comments{Response: chatData}

		log.Println("all data :", comnents)
		c.JSON(http.StatusOK, comnents)
	})

	r.GET("/insertstatus", func(c *gin.Context) {
		chatData := chat.SelectTestData(session, &chatData)
		c.JSON(http.StatusOK, chatData)
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
	err := r.Run(portString)
	if err != nil {
		log.Fatal("gin-gomic run error:", err)
	}

}
