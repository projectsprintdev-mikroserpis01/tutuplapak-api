CREATE TABLE purchase (
    id SERIAL PRIMARY KEY,
    purchased_items JSONB[] NOT NULL,
    sender_name VARCHAR(255) NOT NULL,
    sender_contact_type VARCHAR(50) NOT NULL,
    sender_contact_detail VARCHAR(255) NOT NULL,
    payment_proof_ids VARCHAR(255)[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);