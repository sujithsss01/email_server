package kafkaserver

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sujithsss01/email_server/config"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewKafkaConsumer(topics []string, consumerGroup string, handlerFunc func(*kafka.Message) error) {
	cf := config.GetConfig()
	broker := cf.KafkaConfig.Url
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGINT)
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  broker,
		"group.id":           consumerGroup,
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		log.Printf("Failed to connect to Kafka Consumer")
		return
	}
	defer consumer.Close()
	err = consumer.SubscribeTopics(topics, nil)
	for true {
		select {
		case sig := <-sigchan:
			log.Printf("revieved system signal to stop the consumer: %v", sig)
			return
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				if err := handlerFunc(e); err != nil {
					log.Println("Failed to process the message. closing the consumer")
					return
				}
				if _, err := consumer.CommitMessage(e); err != nil {
					log.Println("Failed to commit the message")
				}
			case kafka.Error:
				log.Printf("Error in consumer: %v", e.Error())
			default:
				log.Printf("dafault case hit")
			}
		}
	}
}
