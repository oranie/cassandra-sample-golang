package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

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

	now := time.Now()

	if err := session.Query(`INSERT INTO chat (name,time,chat_room,comment) VALUES (?,?,?,?)`,
		"oranie", now.UnixNano(), "game_room1", "test comment"+now.String()).Exec(); err != nil {
		log.Fatal(err)
	}

	var name string
	var time int64
	var chat_room string
	var comment string

	if err := session.Query(`SELECT * FROM chat`).Consistency(gocql.One).Scan(&name, &time, &chat_room, &comment); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Chat:", name, time, chat_room, comment)

	// list all tweets
	iter := session.Query(`SELECT name,time,chat_room,comment FROM chat`).Iter()
	for iter.Scan(&name, &time, &chat_room, &comment) {
		fmt.Println("All Chat:", name, time, chat_room, comment)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}
