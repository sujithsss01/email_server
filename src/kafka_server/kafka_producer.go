package kafkaserver

import (
	"email_server/config"
	"errors"
	"log"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	topic        string
	producer     *kafka.Producer
	deliveryChan chan kafka.Event
	mu           sync.Mutex
	closed       bool
}

func (p *KafkaProducer) SendMessage(val []byte, headers map[string]string) error {
	var kafkaHeaders []kafka.Header
	for k, v := range headers {
		kafkaHeaders = append(kafkaHeaders, kafka.Header{
			Key:   k,
			Value: []byte(v),
		})
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return errors.New("This producer is already closed")
	}
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          val,
		Headers:        kafkaHeaders,
	}, p.deliveryChan)
}

func (p *KafkaProducer) Close() {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	p.mu.Unlock()
	rem := p.producer.Flush(15_000)
	if rem > 0 {
		log.Println("warning: there are pending messages in this producer")
	}
	p.producer.Close()
}

func checkForDeliveryReport(deliveryChan chan kafka.Event) {
	for event := range deliveryChan {
		switch e := event.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				log.Printf("failed to deliver the message: %v", e.TopicPartition.Error)
			} else {
				log.Printf("message delivered: %v", e.TopicPartition)
			}

		}
	}
}

func NewKafkaProducer(topic Kafkatopic) (*KafkaProducer, error) {
	cf := config.GetConfig()
	broker := cf.KafkaConfig.Url
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Println("failed to create producer")
		return nil, err
	}
	deliveryChan := make(chan kafka.Event)
	go checkForDeliveryReport(deliveryChan)
	return &KafkaProducer{
		producer:     producer,
		deliveryChan: deliveryChan,
		topic:        string(topic),
	}, nil
}
