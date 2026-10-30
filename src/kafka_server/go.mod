module github.com/sujithsss01/email_server/kafka_server

go 1.27.1

require (
	github.com/confluentinc/confluent-kafka-go/v2 v2.16.0
	github.com/sujithsss01/email_server/config v0.0.0
)

require (
	github.com/caarlos0/env/v11 v11.4.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)

replace github.com/sujithsss01/email_server/config => ../config
