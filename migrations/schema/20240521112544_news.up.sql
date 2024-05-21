CREATE TABLE news (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) UNIQUE,
  content TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
