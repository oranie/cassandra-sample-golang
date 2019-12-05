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
	cluster.CQLVersion = "5.0.1"
	cluster.ProtoVersion = 4
	cluster.Port = 9142
	cluster.DisableInitialHostLookup = true
	cluster.IgnorePeerAddr = true
	/*
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: "cassandra",
			Password: "cassandra",
		}
	*/

	cluster.SslOpts = &gocql.SslOptions{
		CaPath:                 "./AmazonRootCA1.pem",
		EnableHostVerification: false,
	}

	session, error := cluster.CreateSession()
	if error != nil {
		log.Printf("Error: connect cassandra cluster : %v", cluster)
		panic(error.Error())
	}

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
