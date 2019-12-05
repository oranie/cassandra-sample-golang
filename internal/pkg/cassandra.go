package chat

import (
	"github.com/gocql/gocql"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

type Env struct {
	CassandraPort     int    `envconfig:"CASSANDRA_PORT" default:"9042"`
	CassandraUserName string `envconfig:"CASSANDRA_USER" default:"cassandra"`
	CassandraUserPass string `envconfig:"CASSANDRA_PASS" default:"cassandra"`
	CassandraKeyspace string `envconfig:"CASSANDRA_KS" default:"example"`
}

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
	var env Env
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("env set value : ", env)

	cluster := gocql.NewCluster(cassandraCluster)
	cluster.Keyspace = env.CassandraKeyspace
	cluster.Consistency = gocql.Quorum
	cluster.Port = env.CassandraPort
	cluster.DisableInitialHostLookup = true

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: env.CassandraUserName,
		Password: env.CassandraUserPass,
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
