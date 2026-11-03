CREATE TABLE auth (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    auth_type SMALLINT NOT NULL,
    pass_hash VARCHAR(500),

    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
    ,
    CONSTRAINT uq_auth_user_type
        UNIQUE (user_id, auth_type)  
);