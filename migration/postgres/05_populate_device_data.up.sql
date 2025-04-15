-- Xóa dữ liệu cũ nếu cần (chỉ dùng khi dev/test)
TRUNCATE device_history, device_schedule, device_status, devices RESTART IDENTITY;

-- 1. Thêm thiết bị trước và lấy ID tự động
INSERT INTO devices (serial_device, name, type, room) VALUES
    ('SN001', 'Đèn phòng khách', 'LIGHT', 'Phòng khách'),
    ('SN002', 'Quạt phòng ngủ', 'FAN', 'Phòng ngủ'),
    ('SN003', 'Ổ cắm nhà bếp', 'PLUG', 'Nhà bếp');

-- 2. Thêm trạng thái thiết bị (sử dụng ID đúng)
INSERT INTO device_status (device_id, active, schedule_enable, state, updated_at) VALUES
    (1, true, true, '{"power": "ON", "brightness": 70}', NOW()),
    (2, true, false, '{"power": "OFF", "speed": 0}', NOW()),
    (3, false, false, '{"power": "OFF"}', NOW());

-- 3. Thêm lịch trình thiết bị
INSERT INTO device_schedule (device_id, action, scheduled_at, recurring, is_active) VALUES
    (1, 'TURN_OFF', NOW() + INTERVAL '1 hour', 'DAILY', true),
    (2, 'TURN_ON', NOW() + INTERVAL '8 hours', 'ONCE', true);

-- 4. Thêm lịch sử thiết bị
INSERT INTO device_history (device_id, action, state, by, timestamp) VALUES
    (1, 'TURN_ON', '{"power": "ON"}', 'user', NOW() - INTERVAL '2 hours'),
    (1, 'ADJUST', '{"brightness": 70}', 'user', NOW() - INTERVAL '1 hour'),
    (2, 'TURN_OFF', '{"power": "OFF"}', 'auto', NOW() - INTERVAL '30 minutes');