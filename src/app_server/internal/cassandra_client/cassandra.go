package cassandraclient

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/sujithsss01/email_server/config"
)

func GetCassandraSession() *gocql.Session {
	cf := config.GetConfig()
	cluster := gocql.NewCluster(cf.CassandraConfig.Url)
	cluster.Keyspace = cf.CassandraConfig.Keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Errorf("failed to create Cassandra session: %w", err))
	}
	return session
}

func RunMigrations(session *gocql.Session) {
	cf := config.GetConfig()
	driver, err := cassandra.WithInstance(
		session,
		&cassandra.Config{
			KeyspaceName: cf.CassandraConfig.Keyspace,
		},
	)
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"cassandra",
		driver,
	)
	if err != nil {
		panic(err)
	}
	defer m.Close()
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	log.Println("Migrations completed")
}
