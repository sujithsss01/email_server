package main

import (
	"email_server/app_server/internal"
	cassandraclient "email_server/app_server/internal/cassandra_client"
	kafkaserver "email_server/kafka_server"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/gin-gonic/gin"
)

func runAppServer(session *gocql.Session) {
	router := gin.Default()
	app := NewApp(session)
	app.registerRoutes(router)
	router.Run()
}

func main() {
	kafkaserver.InitiateKafkaSetup()
	session := cassandraclient.GetCassandraSession()
	defer session.Close()
	go internal.IntiateIncommingEmailsConsumer(session)
	runAppServer(session)
}
