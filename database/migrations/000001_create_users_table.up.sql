CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  phone VARCHAR(50) NULL UNIQUE,
  bank_account_name VARCHAR(32) NULL,
  bank_account_holder VARCHAR(32) NULL,
  bank_account_number VARCHAR(32) NULL,
  file_id INT NULL,
  file_uri VARCHAR(255) NULL,
  file_thumbnail_uri VARCHAR(255) NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
