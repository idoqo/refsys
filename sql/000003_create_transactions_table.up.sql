CREATE TABLE IF NOT EXISTS transactions
(
    id SERIAL PRIMARY KEY,
    reference VARCHAR(128) UNIQUE NOT NULL,
    amount DECIMAL(19,8),
    sender_id BIGINT,
    recipient_id BIGINT,
    description VARCHAR(128) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);