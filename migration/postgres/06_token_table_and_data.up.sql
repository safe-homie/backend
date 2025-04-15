CREATE TABLE IF NOT EXISTS tokens (
    device_id VARCHAR(30) PRIMARY KEY,
    token VARCHAR(50)
);

INSERT INTO tokens (device_id, token) VALUES ('default-device-id', 'ExponentPushToken[bk_fWXCvBpUU0pnWy7OW5Y]');