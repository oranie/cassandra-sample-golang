package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type Chat struct {
	Name     string
	Time     int64
	Chatroom string
	Comment  string
}

func main() {
	// connect to the cluster
	// My envroiment : local laptop need to connect cassandra cluster with ssh tunnel
	// example ssh ssh.host -L 9042:cassandra.host:9042

	cassandraCluster := os.Getenv("CASSANDRA_CLUSTER")
	port := os.Getenv("PORT")
	if cassandraCluster == "" {
		panic("CassandraCluster endpint is not defind ENV")
	}
	if port == "" {
		panic("App port is not defind ENV")
	}

	cluster := gocql.NewCluster(cassandraCluster)
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	//cluster.CQLVersion = "5.0.1"

	session, error := cluster.CreateSession()
	if error != nil {
		fmt.Println(error)
	}

	defer session.Close()

	createChatTable(session)

	//generate test data
	chatData := generateChatData()

	r := gin.Default()

	r.GET("/run-test", func(c *gin.Context) {
		//insert test data
		insertData(session, &chatData)

		//select insert data
		selectTestData(session, &chatData)

		//select all data at chat table
		allSelectData(session)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/insertstatus", func(c *gin.Context) {
		chatData := selectTestData(session, &chatData)
		json, err := json.Marshal(chatData)
		if err != nil {
			panic(err)
		}
		c.String(http.StatusOK, string(json))
	})
	r.GET("/alldata", func(c *gin.Context) {
		allChatData := allSelectData(session)
		json, err := json.Marshal(allChatData)
		if err != nil {
			panic(err)
		}
		c.String(http.StatusOK, string(json))
	})

	r.Run()
}

// create chat table
func createChatTable(session *gocql.Session) {
	log.Println("Create chat table progress.......")
	if err := session.Query(`CREATE TABLE IF NOT EXISTS chat (
		name text,
		time bigint,
		chat_room text,
		comment text,
		PRIMARY KEY (name, time)) 
		WITH CLUSTERING ORDER BY (time DESC);`).Exec(); err != nil {
		log.Fatal(err)
	}
	log.Println("Create chat table done!")
}

// Insert test data
func insertData(session *gocql.Session, chatData *Chat) {
	log.Println("Insert test data....")
	if err := session.Query(`INSERT INTO chat (name,time,chat_room,comment) VALUES (?,?,?,?)`,
		chatData.Name,
		chatData.Time,
		chatData.Chatroom,
		chatData.Comment).Exec(); err != nil {
		log.Fatal(err)
	}
	log.Println("Insert test data done!")
}

//Get Insert data
func selectTestData(session *gocql.Session, chatData *Chat) Chat {
	log.Println("Select insert test data....")
	var selectChatData Chat
	if err := session.Query(`SELECT name,time,chat_room,comment FROM chat where name = ?`,
		chatData.Name).Consistency(gocql.One).Scan(
		&selectChatData.Name,
		&selectChatData.Time,
		&selectChatData.Chatroom,
		&selectChatData.Comment); err != nil {
		log.Fatal(err)
	}
	log.Println("Insert Data:", selectChatData)
	log.Println("Select insert test data done!")
	return selectChatData
}

// list all chat
func allSelectData(session *gocql.Session) []Chat {
	// list all chat
	log.Println("Select all table data...")
	var ChatData Chat
	selectAllChatData := []Chat{}
	iter := session.Query(`SELECT name,time,chat_room,comment FROM chat`).Iter()
	for iter.Scan(
		&ChatData.Name,
		&ChatData.Time,
		&ChatData.Chatroom,
		&ChatData.Comment) {
		log.Println("All Chat:", ChatData)
		selectAllChatData = append(selectAllChatData, ChatData)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(1000)
	log.Println("Select all table data Done!")
	return selectAllChatData
}

// generate random test data
func generateChatData() Chat {
	now := time.Now()

	chatData := Chat{}
	randomString := random()
	chatData.Name = "oranie-" + randomString
	chatData.Time = now.UnixNano()
	chatData.Chatroom = "game_room-" + randomString
	chatData.Comment = "test comment : " + now.String()

	return chatData
}

//generate random string: 4char
func random() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	rand_string := strconv.FormatUint(n, 36)
	return rand_string[:4]
}
