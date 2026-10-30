package emailmodule

import "github.com/gocql/gocql"

type emailService struct {
	session *gocql.Session
}

func newEmailService(session *gocql.Session) *emailService {
	return &emailService{
		session: session,
	}
}
