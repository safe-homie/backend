CREATE TABLE IF NOT EXISTS sensors(
    id SERIAL PRIMARY KEY,
    type VARCHAR(50),
    location VARCHAR(100) DEFAULT 'n/a',
    unit VARCHAR(20) DEFAULT 'n/a',
    threshold_warning DOUBLE PRECISION DEFAULT 0,
    threshold_danger DOUBLE PRECISION DEFAULT 0
);

CREATE TABLE IF NOT EXISTS sensor_data (
    sensor_id INTEGER,
    time TIMESTAMP NOT NULL DEFAULT NOW(),
    value DOUBLE PRECISION,
    FOREIGN KEY (sensor_id) REFERENCES sensors (id)
);

SELECT create_hypertable('sensor_data', 'time');