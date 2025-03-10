INSERT INTO sensors 
    (type, name, location, unit, threshold_warning, threshold_danger) 
VALUES
    ('humidity', 'Humidity Sensor A', 'living room', '%', 50.0, 80.0),
    ('temperature', 'Temperature Sensor B', 'living room', '°C', 30.0, 50.0),
    ('light', 'Light Sensor C', 'bedroom', 'lx', 100.0, 300.0),
    ('humidity', 'Humidity Sensor D', 'bedroom', '%', 40.0, 70.0);

INSERT INTO sensor_data 
    (sensor_id, time, value) 
VALUES
    (1, NOW() - INTERVAL '5 minutes', 45.5),
    (1, NOW() - INTERVAL '2 minutes', 49.0),
    (2, NOW() - INTERVAL '10 minutes', 25.2),
    (2, NOW() - INTERVAL '3 minutes', 29.5),
    (3, NOW() - INTERVAL '15 minutes', 150.0),
    (3, NOW() - INTERVAL '1 minutes', 280.0),
    (4, NOW() - INTERVAL '7 minutes', 39.0),
    (4, NOW() - INTERVAL '2 minutes', 41.5);