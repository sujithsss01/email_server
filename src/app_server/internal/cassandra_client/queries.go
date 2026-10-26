package cassandraclient

import (
	"fmt"
)

func createKeyspaceQuery(keyspaceName string) string {
	return fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s
			WITH replication = {
				'class': 'SimpleStrategy',
				'replication_factor': 1
			}
	`, keyspaceName)
}

func createEmailByIdTableQuery() string {
	return fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS email_by_id (
			user_id UUID,
			email_id UUID,
			email_type TEXT,
			created_at TIMESTAMP,
			content TEXT,
			PRIMARY KEY((user_id), email_id)
		)
	`)
}
func createEmailByTypeTableQuery() string {
	return fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS email_by_type (
			user_id UUID,
			email_id UUID,
			email_type TEXT,
			created_at TIMESTAMP,
			content TEXT,
			PRIMARY KEY((user_id), email_type, created_at)
		) WITH CLUSTERING ORDER BY (
			email_type ASC,
		    created_at DESC
		)
	`)
}
