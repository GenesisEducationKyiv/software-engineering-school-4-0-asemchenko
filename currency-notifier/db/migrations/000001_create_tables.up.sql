CREATE TABLE IF NOT EXISTS subscriptions (id SERIAL PRIMARY KEY,
                                          email TEXT NOT
                                          NULL UNIQUE,
                                          subscribed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);


CREATE TABLE IF NOT EXISTS exchange_rates (id SERIAL PRIMARY KEY,
                                           rate NUMERIC NOT
                                           NULL,
                                           fetched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);