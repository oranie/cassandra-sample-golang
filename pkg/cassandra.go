package chat

import (
	"fmt"
	"log"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/gocql/gocql"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	AppEnv            string `envconfig:"APP_ENV" default:"test"`
	AppPort           string `envconfig:"APP_PORT" default:"8081"`
	AppEndpoint       string `envconfig:"APP_ENDPOINT" default:"http://127.0.0.1"`
	CassdraEndpoint   string `envconfig:"CASSANDRA_ENDPOINT" default:"127.0.0.1"`
	CassandraPort     int    `envconfig:"CASSANDRA_PORT" default:"9042"`
	CassandraUserName string `envconfig:"CASSANDRA_USER" default:"cassandra"`
	CassandraUserPass string `envconfig:"CASSANDRA_PASS" default:"cassandra"`
	CassandraKeyspace string `envconfig:"CASSANDRA_KS" default:"example"`
}

func GetEnvValue() Env {
	var env Env
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err.Error())
	}
	return env
}

func InitApi() (Env, *gocql.Session, Comment) {
	env := GetEnvValue()

	if env.AppEnv == "prd" || env.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	session, error := CreateCassandraSession()
	if error != nil {
		fmt.Println(error)
	}

	//generate test data
	chatData := GenerateChatData()
	return env, session, chatData
}

func InitGinApp(r *gin.Engine) {
	session, error := CreateCassandraSession()
	if error != nil {
		log.Println(error)
	}

	//if not exist table, create chat table
	CreatePKChatroomTable(session)
}

func CreateSessionConf(env Env) (*gocql.ClusterConfig, Env) {
	cluster := gocql.NewCluster(env.CassdraEndpoint)
	cluster.Keyspace = env.CassandraKeyspace
	cluster.Consistency = gocql.LocalQuorum
	cluster.Port = env.CassandraPort

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: env.CassandraUserName,
		Password: env.CassandraUserPass,
	}

	// Amazon MCS check
	AmazonMCSCheck := regexp.MustCompile(`cassandra.*.amazonaws.com`)
	if AmazonMCSCheck.MatchString(env.CassdraEndpoint) {
		log.Println("This setting connect to Amazon MCS")
		cluster.DisableInitialHostLookup = true
		cluster.SslOpts = &gocql.SslOptions{
			CaPath:                 "./AmazonRootCA1.pem",
			EnableHostVerification: false,
		}
	}
	log.Println("CreateSessionConf is", cluster)
	return cluster, env
}

// connect to the cluster
// My envroiment : local laptop need to connect cassandra cluster with ssh tunnel
// example ssh ssh.host -L 9042:cassandra.host:9042
func CreateCassandraSession() (*gocql.Session, error) {
	env := GetEnvValue()

	cluster, _ := CreateSessionConf(env)
	session, con_error := cluster.CreateSession()
	log.Println("cluster setting is  :", cluster)
	if con_error != nil {
		log.Printf("Error: connect cassandra cluster : %v %v %v", cluster, session, con_error)
		panic(con_error)
	}

	return session, con_error
}

// create chat table
func CreatePKNameTable(session *gocql.Session) {
	log.Println("Create chat(PK:name,time)  table progress.......")
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

func CreatePKChatroomTable(session *gocql.Session) {
	log.Println("Create chat(PK:chatroom,time) table progress.......")
	if err := session.Query(`CREATE TABLE IF NOT EXISTS chat (
		name text,
		time bigint,
		chatroom text,
		comment text,
		PRIMARY KEY (chatroom, time)) 
		WITH CLUSTERING ORDER BY (time DESC);`).Exec(); err != nil {
		log.Fatal(err)
	}
	log.Println("Create chat table done!")
}
