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
);