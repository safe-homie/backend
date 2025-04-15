CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    serial_device VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(50),
    type VARCHAR(50),
    room VARCHAR(100)
);
CREATE TABLE IF NOT EXISTS device_status (
    device_id INT PRIMARY KEY REFERENCES devices(id) ON DELETE CASCADE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    schedule_enable BOOLEAN NOT NULL DEFAULT FALSE,
    state JSONB NOT NULL, -- power ON/OFF, level MIN/MEDIUM/MAX
    updated_at TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS device_history (
    id SERIAL PRIMARY KEY,
    device_id INT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    action TEXT NOT NULL,
    state JSONB NOT NULL,
    by TEXT NOT NULL, -- "user" or "scheduling:%ID"
    timestamp TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS device_schedule (
    id SERIAL PRIMARY KEY,
    device_id INT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    action TEXT NOT NULL CHECK (action IN ('TURN_ON', 'TURN_OFF', 'SET_LEVEL')),
    scheduled_at TIMESTAMP(0) NOT NULL,
    recurring TEXT NOT NULL CHECK (recurring IN ('DAILY', 'WEEKLY', 'ONCE')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS  idx_device_history_device_id ON device_history(device_id);
CREATE INDEX IF NOT EXISTS  idx_device_schedule_device_id ON device_schedule(device_id);
CREATE INDEX IF NOT EXISTS  idx_device_status_device_id ON device_status(device_id);
