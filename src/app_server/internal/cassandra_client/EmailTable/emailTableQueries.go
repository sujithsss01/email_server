package emailtablequeries

import (
	"time"

	"github.com/gocql/gocql"
)

func SaveEmail(
	session *gocql.Session,
	user_id gocql.UUID,
	email_id gocql.UUID,
	email_type string,
	created_at time.Time,
	content string,
) error {
	batch := session.NewBatch(gocql.LoggedBatch)
	batch.Query(`
		INSERT INTO email_by_id
			(user_id, email_id, email_type, created_at, content)
		VALUES ( ?, ?, ?, ?, ? )
	`, user_id, email_id, email_type, created_at, content)
	batch.Query(`
		INSERT INTO email_by_type
			(user_id, email_id, email_type, created_at, content)
		VALUES ( ?, ?, ?, ?, ? )
	`, user_id, email_id, email_type, created_at, content)
	return session.ExecuteBatch(batch)

}
