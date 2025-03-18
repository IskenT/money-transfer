-- +migrate Up
CREATE SCHEMA IF NOT EXISTS money_transfer;

-- Users table
CREATE TABLE money_transfer.users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Transactions table
CREATE TABLE money_transfer.transactions (
    id BIGSERIAL PRIMARY KEY,
    stan VARCHAR(50) NOT NULL,
    amount BIGINT NOT NULL,
    state VARCHAR(20) NOT NULL,
    transaction_type VARCHAR(20) NOT NULL,
    payment_source VARCHAR(20) NOT NULL,
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_transactions_stan ON money_transfer.transactions(stan);

-- Transfers table
CREATE TABLE money_transfer.transfers (
    id SERIAL PRIMARY KEY,
    transfer_code VARCHAR(50) UNIQUE NOT NULL,
    from_user_id INT NOT NULL REFERENCES money_transfer.users(id),
    to_user_id INT NOT NULL REFERENCES money_transfer.users(id),
    amount BIGINT NOT NULL CHECK (amount > 0),
    state VARCHAR(20) NOT NULL,
    debit_tx_id BIGINT REFERENCES money_transfer.transactions(id),
    credit_tx_id BIGINT REFERENCES money_transfer.transactions(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX idx_transfers_users ON money_transfer.transfers(from_user_id, to_user_id);

-- Outbox 
CREATE TABLE money_transfer.outbox_events (
    id BIGSERIAL PRIMARY KEY,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id VARCHAR(50) NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX idx_outbox_unprocessed ON money_transfer.outbox_events(processed_at) WHERE processed_at IS NULL;

-- Initial seed data
INSERT INTO money_transfer.users (id, name, balance) VALUES 
    (1, 'Mark', 10000),
    (2, 'Jane', 5000),
    (3, 'Adam', 0);

-- Reset 
SELECT setval('money_transfer.users_id_seq', (SELECT MAX(id) FROM money_transfer.users));

-- +migrate Down
DROP TABLE IF EXISTS money_transfer.outbox_events;
DROP TABLE IF EXISTS money_transfer.transfers;
DROP TABLE IF EXISTS money_transfer.transactions;
DROP TABLE IF EXISTS money_transfer.users;
DROP SCHEMA IF EXISTS money_transfer;