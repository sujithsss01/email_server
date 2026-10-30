CREATE TABLE IF NOT EXISTS email_by_id (
    user_id UUID,
    email_id UUID,
    email_type TEXT,
    created_at TIMESTAMP,
    content TEXT,
    PRIMARY KEY((user_id), email_id)
);