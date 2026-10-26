package main

import (
	"encoding/json"
	"log"

	"github.com/sujithsss01/email_server/kafka_server"
	"github.com/sujithsss01/email_server/smtp_server/dto"
)

func processEmail(producer *kafkaserver.KafkaProducer, to, from, subject, message string) error {
	email := dto.Email{
		To:      to,
		From:    from,
		Subject: subject,
		Message: message,
	}
	jsonString, err := json.Marshal(email)
	if err != nil {
		log.Println("Failed to process to json")
		return err
	}
	return producer.SendMessage(jsonString, make(map[string]string))
}
