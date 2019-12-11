package main

import (
	"log"
	"regexp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/oranie/cassandra-sample-golang/pkg"
)

func main() {
	r := gin.Default()
	env := chat.GetEnvValue()
	chat.InitGinApp(r)

	ApiEndpoint := "http://localhost:" + env.AppPort + "/"
	// local env check
	localCheck := regexp.MustCompile(`localhost|127.0.0.1`)
	if !localCheck.MatchString(env.AppEndpoint) {
		ApiEndpoint = "https://" + env.AppEndpoint + "/"
	}
	log.Println("App Endpoint : ", ApiEndpoint)

	r.Use(cors.Default())

	//Server status check endpoint.
	r.GET("/", chat.GetServerStatus)

	//Insert test data and select.
	r.GET("/run-test", chat.GetRunTestResult)

	//Return chat html page
	r.GET("/chat", chat.GetHTMLPage)

	//Insert Chat comment data
	r.POST("/chat/comments/add", chat.PostInsertChatData)

	//Get Latest Chat comment data
	r.GET("/chat/comments/latest", chat.GetLatestChatData)

	//Get All Chat comment data
	r.GET("/chat/comments/all", chat.GetAllChatData)

	portString := ":" + env.AppPort
	err := r.Run(portString)
	if err != nil {
		log.Fatal("gin-gomic run error:", err)
	}

}
