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
	// local laptop need to connect ssh tunnel
	// example ssh ssh.host -L 9042:cassandra.host:9042

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	cluster.CQLVersion = "5.0.1"
	session, error := cluster.CreateSession()
	if error != nil {
		fmt.Println(error)
	}

	defer session.Close()

	// create table
	println("Create chat table progress.......")
	if err := session.Query(`CREATE TABLE IF NOT EXISTS chat (
		name text,
		time bigint,
		chat_room text,
		comment text,
		PRIMARY KEY (name, time)) 
		WITH CLUSTERING ORDER BY (time DESC);`).Exec(); err != nil {
		log.Fatal(err)
	}

	println("Create chat table done!")

	var name string
	var time int64
	var chat_room string
	var comment string

	chatData := generateChatData()

	if err := session.Query(`INSERT INTO chat (name,time,chat_room,comment) VALUES (?,?,?,?)`,
		chatData.Name, chatData.Time, chatData.Chat_room, chatData.Comment).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := session.Query(`SELECT * FROM chat where name = ?`,
		chatData.Name).Consistency(gocql.One).Scan(&name, &time, &chat_room, &comment); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Insert Chat:", name, time, chat_room, comment)

	// list all chat
	iter := session.Query(`SELECT name,time,chat_room,comment FROM chat`).Iter()
	for iter.Scan(&name, &time, &chat_room, &comment) {
		fmt.Println("All Chat:", name, time, chat_room, comment)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}

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

func random() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	rand_string := strconv.FormatUint(n, 36)
	return rand_string[:4]
}
