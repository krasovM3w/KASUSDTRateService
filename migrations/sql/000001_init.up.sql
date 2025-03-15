CREATE TABLE IF NOT EXISTS rates (
                                     id SERIAL PRIMARY KEY,
                                     base_currency VARCHAR(10) NOT NULL,
    target_currency VARCHAR(10) NOT NULL,
    rate DECIMAL(15, 6) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                            );

CREATE INDEX idx_rates_timestamp ON rates(timestamp);