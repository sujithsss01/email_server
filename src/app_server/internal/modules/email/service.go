package emailmodule

import gocql "github.com/apache/cassandra-gocql-driver/v2"

type emailService struct {
	session *gocql.Session
}

func newEmailService(session *gocql.Session) *emailService {
	return &emailService{
		session: session,
	}
}
