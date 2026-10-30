package internal

import (
	"encoding/json"
	emailtablequeries "github.com/sujithsss01/email_server/app_server/internal/cassandra_client/emailTable"
	kafkaserver "github.com/sujithsss01/email_server/kafka_server"
	"log"

	emaildto "github.com/sujithsss01/email_server/smtp_server/dto"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gocql/gocql"
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
