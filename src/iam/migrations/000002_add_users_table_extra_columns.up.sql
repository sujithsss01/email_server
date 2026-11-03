ALTER TABLE users
ADD COLUMN phone_number VARCHAR(30) NOT NULL UNIQUE DEFAULT '',
ADD COLUMN active boolean NOT NULL DEFAULT FALSE,
ADD CONSTRAINT check_phone_number
    CHECK(REGEXP_LIKE(phone_number, '^\+[1-9][1-9]{6,14}$'));