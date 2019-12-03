package chat

import (
	"github.com/gocql/gocql"
	"log"
	"os"
)

func CreateCassandraSession() (*gocql.Session, error) {
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
	return session, error
}

// create chat table
func CreateChatTable(session *gocql.Session) {
	log.Println("Create chat table progress.......")
	if err := session.Query(`CREATE TABLE IF NOT EXISTS chat (
		name text,
		time bigint,
		chatroom text,
		comment text,
		PRIMARY KEY (name, time)) 
		WITH CLUSTERING ORDER BY (time DESC);`).Exec(); err != nil {
		log.Fatal(err)
	}
	log.Println("Create chat table done!")
}
