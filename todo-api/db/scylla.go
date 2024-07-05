package db

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func InitDB() {
	cluster := gocql.NewCluster(os.Getenv("SCYLLA_HOST"))
	cluster.Port = getEnvAsInt("SCYLLA_PORT", 9042)
	cluster.Keyspace = os.Getenv("SCYLLA_KEYSPACE")
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: os.Getenv("SCYLLA_USERNAME"),
		Password: os.Getenv("SCYLLA_PASSWORD"),
	}
	cluster.ProtoVersion = 4
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = time.Second * 10
	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("Unable to connect to ScyllaDB: ", err)
	} else {
		log.Println("Successfully connected to ScyllaDB")
	}
}

func getEnvAsInt(name string, defaultVal int) int {
	if value, exists := os.LookupEnv(name); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultVal
}
