CREATE TABLE simpleStorageContract (
  id SERIAL PRIMARY KEY,
  value BIGINT NOT NULL,
  timestamp TIMESTAMP DEFAULT NOW()
);
