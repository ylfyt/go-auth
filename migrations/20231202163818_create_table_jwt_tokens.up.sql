CREATE TABLE IF NOT EXISTS jwt_tokens (
  id BIGINT PRIMARY KEY,
  created_at timestamp NOT NULL,
  updated_at timestamp NULL,
  user_id INTEGER NOT NULL
);