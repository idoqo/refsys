CREATE TYPE payout_type AS ENUM('signups', 'transactions');
CREATE TYPE payout_status AS ENUM('pending', 'paid', 'error');

CREATE TABLE IF NOT EXISTS payouts (
    id SERIAL PRIMARY KEY,
    user_id BIGINT,
    checkpoint_id BIGINT,
    status payout_status DEFAULT 'pending',
    status_description VARCHAR(128) DEFAULT NULL,
    activity_type payout_type,
    amount DECIMAL(19, 4)
);