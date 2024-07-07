CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    email VARCHAR(200) UNIQUE,
    created_at TIMESTAMP default (now() at time zone 'utc')
    );