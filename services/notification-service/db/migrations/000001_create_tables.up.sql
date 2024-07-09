CREATE TABLE IF NOT EXISTS subscriptions
(
    id            SERIAL PRIMARY KEY,
    email         TEXT NOT NULL UNIQUE,
    subscribed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);