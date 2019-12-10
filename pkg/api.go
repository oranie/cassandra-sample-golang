package chat

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

func GetServerStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "server status id good",
	})
}

func GetRunTestResult(c *gin.Context) {
	_, session, chatData := InitApi()
	if &chatData != nil {
		chatData = GenerateChatData()
	}
	//insert test data
	InsertData(session, &chatData)
	//select insert data
	SelectTestData(session, &chatData)
	//select all data at chat table
	result := AllSelectData(session)
	c.String(http.StatusOK, "Test done.", result)
}

func GetHTMLPage(c *gin.Context) {
	env, _, _ := InitApi()
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
}

func PostInsertChatData(c *gin.Context) {
	_, session, _ := InitApi()
	defer session.Close()

	var json Comment
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postData := Comment{
		Name:     json.Name,
		Time:     time.Now().UnixNano(),
		Chatroom: "game_room-oranie",
		Comment:  json.Comment,
	}

	fmt.Printf("%v", json)
	resp := InsertData(session, &postData)

	c.JSON(http.StatusOK, resp)
}

func GetLatestChatData(c *gin.Context) {
	_, session, _ := InitApi()
	chatroom := "game_room-oranie"
	chatData := ChatroomLatestData(session, chatroom)
	comnents := Comments{Response: chatData}

	log.Println("latest data :", comnents)
	c.JSON(http.StatusOK, comnents)
}

func GetAllChatData(c *gin.Context) {
	_, session, _ := InitApi()
	chatroom := "game_room-oranie"
	chatData := ChatroomAllData(session, chatroom)
	comnents := Comments{Response: chatData}

	log.Println("all data :", comnents)
	c.JSON(http.StatusOK, comnents)
}
