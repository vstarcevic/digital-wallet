CREATE TABLE "balance" (
    userId INT PRIMARY KEY,
    balance DECIMAL(14,2) not null default 0,
    created_at TIMESTAMP default (now() at time zone 'utc')
    );

CREATE TABLE "transaction" (
    transactionId SERIAL PRIMARY KEY,
    userId INT NOT NULL,
    amount DECIMAL(14,2) NOT NULL,
    created_at TIMESTAMP default (now() at time zone 'utc')
    );