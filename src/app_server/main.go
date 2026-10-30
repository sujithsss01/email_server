package main

import (
	"github.com/sujithsss01/email_server/app_server/internal"
	cassandraclient "github.com/sujithsss01/email_server/app_server/internal/cassandra_client"
	kafkaserver "github.com/sujithsss01/email_server/kafka_server"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func runAppServer(session *gocql.Session) {
	router := gin.Default()
	app := NewApp(session)
	app.registerRoutes(router)
	router.Run()
}

func main() {
	session := cassandraclient.GetCassandraSession()
	defer session.Close()
	cassandraclient.RunMigrations(session)
	kafkaserver.InitiateKafkaSetup()
	go internal.IntiateIncommingEmailsConsumer(session)
	runAppServer(session)
}
