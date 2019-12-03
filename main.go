package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gocql/gocql"
)

type Chat struct {
	Name      string
	Time      int64
	Chat_room string
	Comment   string
}

func main() {
	// connect to the cluster
	// My envroiment : local laptop need to connect cassandra cluster with ssh tunnel
	// example ssh ssh.host -L 9042:cassandra.host:9042

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	//cluster.CQLVersion = "5.0.1"

	session, error := cluster.CreateSession()
	if error != nil {
		fmt.Println(error)
	}

	defer session.Close()

	// create table
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

	//Insert test data
	chatData := generateChatData()

	insertTestData(session, &chatData)
	selectTestData(session, &chatData)
	allSelectTestData(session)

}

// Insert test data
func insertTestData(session *gocql.Session, chatData *Chat) {
	log.Println("Insert test data....")
	if err := session.Query(`INSERT INTO chat (name,time,chat_room,comment) VALUES (?,?,?,?)`,
		chatData.Name,
		chatData.Time,
		chatData.Chat_room,
		chatData.Comment).Exec(); err != nil {
		log.Fatal(err)
	}
	log.Println("Insert test data done!")
}

//Get Insert data
func selectTestData(session *gocql.Session, chatData *Chat) {
	log.Println("Select insert test data....")
	var selectChatData Chat
	if err := session.Query(`SELECT name,time,chat_room,comment FROM chat where name = ?`,
		chatData.Name).Consistency(gocql.One).Scan(
		&selectChatData.Name,
		&selectChatData.Time,
		&selectChatData.Chat_room,
		&selectChatData.Comment); err != nil {
		log.Fatal(err)
	}
	log.Println("Insert Data:", selectChatData)
	log.Println("Select insert test data done!")
}

// list all chat
func allSelectTestData(session *gocql.Session) {
	// list all chat
	log.Println("Select all table data...")
	var selectAllChatData Chat
	iter := session.Query(`SELECT name,time,chat_room,comment FROM chat`).Iter()
	for iter.Scan(
		&selectAllChatData.Name,
		&selectAllChatData.Time,
		&selectAllChatData.Chat_room,
		&selectAllChatData.Comment) {
		fmt.Println("All Chat:", selectAllChatData)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
	log.Println("Select all table data Done!")
}

// generate random test data
func generateChatData() Chat {
	now := time.Now()

	chatData := Chat{}
	randomString := random()
	chatData.Name = "oranie-" + randomString
	chatData.Time = now.UnixNano()
	chatData.Chat_room = "game_room-" + randomString
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
