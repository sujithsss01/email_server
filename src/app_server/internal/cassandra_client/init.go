package cassandraclient

import (
	"github.com/sujithsss01/email_server/config"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func GetCassandraSession() *gocql.Session {
	cf := config.GetConfig()
	cluster := gocql.NewCluster(cf.CassandraConfig.Url)
	session, err := cluster.CreateSession()
	if err != nil {
		panic("Failed to create cassandra session")
	}
	err = session.Query(createKeyspaceQuery(cf.CassandraConfig.Keyspace)).Exec()
	if err != nil {
		panic("Failed at keyspace creation")
	}
	session.Close()
	cluster.Keyspace = cf.CassandraConfig.Keyspace
	session, err = cluster.CreateSession()
	if err != nil {
		panic("failed to connect to the cassandra keyspace")
	}
	err = session.Query(createEmailByIdTableQuery()).Exec()
	if err != nil {
		panic("failed to create email table 1")
	}
	err = session.Query(createEmailByTypeTableQuery()).Exec()
	if err != nil {
		panic("failed to create email table 2")
	}
	return session
}
