INSERT INTO transactions(reference, amount, sender_id, recipient_id) VALUES ('seedertransaction', 15000, 0, 1);

INSERT INTO wallets(user_id, transaction_type, transaction_reference, amount, closing_balance) VALUES (1, 'credit', 'seedertransaction', 15000, 15000);
