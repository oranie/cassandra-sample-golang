package chat

import (
	"github.com/gocql/gocql"
	"log"
	"os"
)

// connect to the cluster
// My envroiment : local laptop need to connect cassandra cluster with ssh tunnel
// example ssh ssh.host -L 9042:cassandra.host:9042
func CreateCassandraSession() (*gocql.Session, error) {

	cassandraCluster := os.Getenv("CASSANDRA_ENDPOINT")
	if cassandraCluster == "" {
		panic("Cassandra endpoint is not defind ENV")
	}
	cluster := CreateSessionConf(cassandraCluster)

	session, error := cluster.CreateSession()
	if error != nil {
		log.Printf("Error: connect cassandra cluster : %v", cluster)
		panic(error.Error())
	}

	return session, error
}

func CreateSessionConf(cassandraCluster string) *gocql.ClusterConfig {
	cassandraUserName := os.Getenv("CASSANDRA_USER")
	cassandraUserPass := os.Getenv("CASSANDRA_PASS")
	cassandraKeyspace := os.Getenv("CASSANDRA_KS")

	cluster := gocql.NewCluster(cassandraCluster)
	cluster.Keyspace = cassandraKeyspace
	cluster.Consistency = gocql.Quorum
	cluster.Port = 9142
	cluster.DisableInitialHostLookup = true

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cassandraUserName,
		Password: cassandraUserPass,
	}

	cluster.SslOpts = &gocql.SslOptions{
		CaPath:                 "./AmazonRootCA1.pem",
		EnableHostVerification: false,
	}

	return cluster
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
