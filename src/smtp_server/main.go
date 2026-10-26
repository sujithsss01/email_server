package main

import (
	"io"
	"log"
	"net/mail"
	"time"

	"github.com/sujithsss01/email_server/kafka_server"
	"github.com/emersion/go-smtp"
)

type Backend struct {
	producer *kafkaserver.KafkaProducer
}

func (be *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		producer: be.producer,
	}, nil
}

type Session struct {
	from     string
	to       []string
	producer *kafkaserver.KafkaProducer
}

func (s *Session) Data(r io.Reader) error {
	msg, err := mail.ReadMessage(r)
	if err != nil {
		log.Println("invalid message format")
		return err
	}
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		log.Println("invalid message body")
		return err
	}
	to, from := msg.Header.Get("To"), s.from
	subject, message := msg.Header.Get("Subject"), string(body)
	processEmail(s.producer, to, from, subject, message)
	return nil
}
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.to = append(s.to, to)
	return nil
}
func (s *Session) Reset() {
	s.from = ""
	s.to = nil
}
func (s *Session) Logout() error {
	return nil
}
func main() {
	producer, _ := kafkaserver.NewKafkaProducer(kafkaserver.IncommingEmailsTopic)
	defer producer.Close()
	be := &Backend{
		producer: producer,
	}
	server := smtp.NewServer(be)
	server.Addr = "localhost:1025"
	server.Domain = "localhost"
	server.WriteTimeout = 10 * time.Second
	server.ReadTimeout = 10 * time.Second
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true

	log.Println("Starting server at", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
