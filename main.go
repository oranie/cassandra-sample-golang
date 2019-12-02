package main

import (
	"fmt"
	"log"

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
		timestamp text,
		chat_room text,
		comment text,
		PRIMARY KEY (name, timestamp)) 
		WITH CLUSTERING ORDER BY (timestamp DESC);`).Exec(); err != nil {
		log.Fatal(err)
	}

	println("Create chat table done!")

	if err := session.Query(`INSERT INTO chat (name,timestamp,chat_room,comment) VALUES (?,?,?,?)`,
		"oranie", "001", "game_room1", "test comment").Exec(); err != nil {
		log.Fatal(err)
	}

	// insert a tweet
	if err := session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), "hello world").Exec(); err != nil {
		log.Fatal(err)
	}

	var id gocql.UUID
	var text string

	/* Search for a specific set of records whose 'timeline' column matches
	 * the value 'me'. The secondary index that we created earlier will be
	 * used for optimizing the search */
	if err := session.Query(`SELECT id, text FROM tweet WHERE timeline = ? LIMIT 1`,
		"me").Consistency(gocql.One).Scan(&id, &text); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tweet:", id, text)

	// list all tweets
	iter := session.Query(`SELECT id, text FROM tweet WHERE timeline = ?`, "me").Iter()
	for iter.Scan(&id, &text) {
		fmt.Println("Tweet:", id, text)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}
