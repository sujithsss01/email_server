package kafkaserver

import (
	"context"
	"github.com/sujithsss01/email_server/config"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func createIncommingEmailsTopic() *kafka.TopicSpecification {
	return &kafka.TopicSpecification{
		Topic:             string(IncommingEmailsTopic),
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
}

func InitiateKafkaSetup() {
	cf := config.GetConfig()
	broker := cf.KafkaConfig.Url
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Println("failed to create admin client")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	maxDuration, _ := time.ParseDuration("60s")
	results, err := admin.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{
			*createIncommingEmailsTopic(),
		},
		kafka.SetAdminOperationTimeout(maxDuration),
	)
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			log.Printf("Failed to create topic: %v\n", result.Error)
			continue
		}
		log.Printf("Topic created: %s\n", result.Topic)
	}

}
