package internal

import (
	emailtablequeries "email_server/app_server/internal/cassandra_client/EmailTable"
	kafkaserver "email_server/kafka_server"
	"encoding/json"
	"log"

	emaildto "email_server/smtp_server/dto"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

func IntiateIncommingEmailsConsumer(session *gocql.Session) {
	topics := []string{string(kafkaserver.IncommingEmailsTopic)}
	consumer_group := "email_storage"
	handleIncommingEmails := func(message *kafka.Message) error {
		user_id, err := uuid.NewV7()
		if err != nil {
			log.Println("failed to create user_id")
			return err
		}
		mail_id, err := uuid.NewV7()
		if err != nil {
			log.Println("failed to create mail_id")
			return err
		}
		timestamp := message.Timestamp
		var email emaildto.Email
		err = json.Unmarshal(message.Value, &email)
		if err != nil {
			log.Println("Error trying to unpack the json")
			return err
		}
		err = emailtablequeries.SaveEmail(
			session,
			gocql.UUID(user_id),
			gocql.UUID(mail_id),
			"inbox",
			timestamp,
			string(message.Value),
		)
		if err != nil {
			log.Println("Error trying to store the email")
			return err
		}
		return nil
	}
	kafkaserver.NewKafkaConsumer(topics, consumer_group, handleIncommingEmails)
}
