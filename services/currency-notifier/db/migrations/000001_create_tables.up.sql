CREATE TABLE IF NOT EXISTS users (id UUID PRIMARY KEY,
                                  email TEXT NOT NULL UNIQUE,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);


CREATE TABLE IF NOT EXISTS exchange_rates (id SERIAL PRIMARY KEY,
                                           rate NUMERIC NOT NULL,
                                           fetched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);