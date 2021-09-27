CREATE TYPE TXN_TYPE AS ENUM('credit', 'debit');

CREATE TABLE IF NOT EXISTS wallets
(
    id SERIAL PRIMARY KEY,
    user_id BIGINT,
    transaction_type TXN_TYPE,
    transaction_reference VARCHAR(128),
    amount DECIMAL(19,8),
    closing_balance DECIMAL(19,8),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
