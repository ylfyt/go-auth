CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at timestamp NOT NULL,
  updated_at timestamp NULL,
  username varchar NOT NULL,
  "password" varchar NOT NULL
);