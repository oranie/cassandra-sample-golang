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

// connect to the cluster
// My envroiment : local laptop need to connect cassandra cluster with ssh tunnel
// example ssh ssh.host -L 9042:cassandra.host:9042
func CreateCassandraSession() (*gocql.Session, error) {
	env := GetEnvValue()
	log.Println("create session env :", env)
	cluster, _ := CreateSessionConf(env)
	log.Println("cluster.DisableInitialHostLookup", cluster.DisableInitialHostLookup)
	session, error := cluster.CreateSession()
	if error != nil {
		log.Printf("Error: connect cassandra cluster : %v %v %v", cluster, session, error)
		panic(error)
	}

	return session, error
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
		fmt.Println(error)
	}

	//check and create chat table
	CreateChatTable(session)
}

func CreateSessionConf(env Env) (*gocql.ClusterConfig, Env) {
	cluster := gocql.NewCluster(env.CassdraEndpoint)
	cluster.Keyspace = env.CassandraKeyspace
	cluster.Consistency = gocql.Quorum
	cluster.Port = env.CassandraPort

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: env.CassandraUserName,
		Password: env.CassandraUserPass,
	}

	// local env check
	localCheck := regexp.MustCompile(`localhost|127.0.0.1`)
	if !localCheck.MatchString(env.CassdraEndpoint) {
		log.Println("not localhost|127.0.0.1 envroiment")
		cluster.DisableInitialHostLookup = true
		cluster.SslOpts = &gocql.SslOptions{
			CaPath:                 "./AmazonRootCA1.pem",
			EnableHostVerification: false,
		}
	}

	return cluster, env
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
