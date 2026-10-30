module github.com/sujithsss01/email_server/smtp_server

go 1.27.1

require (
	github.com/emersion/go-smtp v0.25.0
	github.com/sujithsss01/email_server/kafka_server v0.0.0
)

require (
	github.com/caarlos0/env/v11 v11.4.1 // indirect
	github.com/confluentinc/confluent-kafka-go/v2 v2.16.0 // indirect
	github.com/emersion/go-sasl v0.0.0-20241020182733-b788ff22d5a6 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/sujithsss01/email_server/config v0.0.0 // indirect
)

replace github.com/sujithsss01/email_server/config => ../config

replace github.com/sujithsss01/email_server/kafka_server => ../kafka_server
