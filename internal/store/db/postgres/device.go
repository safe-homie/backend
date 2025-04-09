package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/safe-homie/backend/internal/store"
)

func (p *_postgres) ListDevices(find *store.FindDevice) ([]*store.Device, error) {
	where, args := []string{}, []any{}
	if v := find.Location; v != nil {
		where, args = append(where, "location = $1"), append(args, *v)
	}
	stmt := `SELECT 
                id,
                location,
                name,
                type,
                attributes->>'power' AS power,
                (attributes->>'level')::INTEGER AS level,
                created_at,
                updated_at
             FROM devices`
	if len(where) > 0 {
		stmt += " WHERE " + strings.Join(where, " AND ")
	}
	if v := find.Limit; v != nil {
		stmt += fmt.Sprintf(" LIMIT %d", *v)
	}
	rows, err := p.db.Query(context.Background(), stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]*store.Device, 0)
	for rows.Next() {
		var device store.Device
		if err := rows.Scan(
			&device.ID,
			&device.Location,
			&device.Name,
			&device.Type,
			&device.Power,
			&device.Level,
			&device.CreateAt,
			&device.UpdateAt,
		); err != nil {
			return nil, err
		}
		list = append(list, &device)
	}
	return list, nil
}
func (p *_postgres) GetDeviceByID(find *store.FindDevice) (*store.Device, error) {
	where, args := []string{}, []any{}
	if v := find.ID; v != nil {
		where, args = append(where, "id = $1"), append(args, *v)
	}
	stmt := `SELECT 
                id,
                location,
                name,
                type,
                attributes->>'power' AS power,
                (attributes->>'level')::INTEGER AS level,
                created_at,
                updated_at
             FROM devices`
	if len(where) > 0 {
		stmt += " WHERE " + strings.Join(where, " AND ")
	}
	if v := find.Limit; v != nil {
		stmt += fmt.Sprintf(" LIMIT %d", *v)
	}
	var device store.Device
	if err := p.db.QueryRow(context.Background(), stmt, args...).Scan(
		&device.ID,
		&device.Location,
		&device.Name,
		&device.Type,
		&device.Power,
		&device.Level,
		&device.CreateAt,
		&device.UpdateAt,
	); err != nil {
		return nil, err
	}
	return &device, nil
}
func (p *_postgres) ListDeviceHistory(find *store.FindDeviceHistory) ([]*store.DeviceHistory, error) {
	// Xây dựng điều kiện WHERE và tham số
	where, args := []string{}, []any{}
	argIndex := 1 // Bắt đầu từ $1

	if v := find.DeviceID; v != nil {
		where = append(where, fmt.Sprintf("device_id = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := find.StartTime; v != nil {
		where = append(where, fmt.Sprintf("timestamp >= $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := find.EndTime; v != nil {
		where = append(where, fmt.Sprintf("timestamp <= $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}

	// Xây dựng câu truy vấn
	stmt := `SELECT
				id,
				device_id,
				action,
				timestamp
			 FROM device_history`
	if len(where) > 0 {
		stmt += " WHERE " + strings.Join(where, " AND ")
	}
	stmt += " ORDER BY timestamp DESC"
	if v := find.Limit; v != nil && *v > 0 {
		stmt += fmt.Sprintf(" LIMIT %d", *v)
	}

	// Thực thi truy vấn
	rows, err := p.db.Query(context.Background(), stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("lỗi truy vấn cơ sở dữ liệu: %v", err)
	}
	defer rows.Close()

	// Duyệt và lấy dữ liệu
	list := make([]*store.DeviceHistory, 0)
	for rows.Next() {
		var history store.DeviceHistory
		if err := rows.Scan(
			&history.ID,
			&history.DeviceID,
			&history.Action,
			&history.Timestamp,
		); err != nil {
			return nil, fmt.Errorf("lỗi quét dữ liệu: %v", err)
		}
		list = append(list, &history)
	}

	// Kiểm tra lỗi sau khi duyệt rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi duyệt hàng: %v", err)
	}

	return list, nil
}

func (p *_postgres) ListDeviceSchedule(find *store.FindDeviceSchedule) ([]*store.DeviceSchedule, error) {
	where, args := []string{}, []any{}
	argIndex := 1

	if v := find.DeviceID; v != nil {
		where = append(where, fmt.Sprintf("device_id = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	stmt := `SELECT
				id,
				device_id,
				action,
				CAST(start_date AS TIMESTAMPTZ) + time AS scheduled_at,
				repeat AS recurring,
				is_active,
				created_at
			 FROM device_schedules`
	if len(where) > 0 {
		stmt += " WHERE " + strings.Join(where, " AND ")
	}
	stmt += " ORDER BY scheduled_at ASC"
	if v := find.Limit; v != nil && *v > 0 {
		stmt += fmt.Sprintf(" LIMIT %d", *v)
	}

	rows, err := p.db.Query(context.Background(), stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*store.DeviceSchedule, 0)
	for rows.Next() {
		var schedule store.DeviceSchedule
		if err := rows.Scan(
			&schedule.ID,
			&schedule.DeviceID,
			&schedule.Action,
			&schedule.ScheduledAt,
			&schedule.Recurring,
			&schedule.IsActive,
			&schedule.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, &schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (p *_postgres) CreateDevice(create *store.Device) (*store.Device, error) {
	stmt := `INSERT INTO devices (location, name, type, attributes, created_at, updated_at)
	VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	RETURNING id, location, name, type, attributes->>'power' AS power, attributes->>'level' AS level, created_at, updated_at`

	levelStr := create.Level
	_, err := strconv.Atoi(levelStr)
	if err != nil {
		return nil, fmt.Errorf("invalid level format: %v", err)
	}

	attributes := fmt.Sprintf(`{"power": "%s", "level": "%s"}`, create.Power, levelStr)

	err = p.db.QueryRow(context.Background(), stmt, create.Location, create.Name, create.Type, attributes).
		Scan(&create.ID, &create.Location, &create.Name, &create.Type, &create.Power, &create.Level, &create.CreateAt, &create.UpdateAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}
	return create, nil
}
func (p *_postgres) UpdateDevice(update *store.UpdateDevice) (*store.UpdateDevice, error) {
	stmt := `
	UPDATE devices SET
		location = COALESCE($1, location),
		name     = COALESCE($2, name),
		type     = COALESCE($3, type),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $4
	RETURNING 
		id, location, name, type, 
		attributes->>'power' AS power, 
		attributes->>'level' AS level,
		created_at, updated_at;
`

	row := p.db.QueryRow(
		context.Background(), stmt,
		update.Location, update.Name, update.Type,
		update.ID,
	)

	device := &store.UpdateDevice{}
	err := row.Scan(
		&device.ID, &device.Location, &device.Name, &device.Type,
		&device.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return device, nil
}

// func (p *_postgres) InsertDeviceHistory(insert *store.DeviceHistory) (*store.DeviceHistory, error) {
// 	fields := []string{"device_id", "action", "timestamp", "user_id"}
// 	args := []any{insert.DeviceID, insert.Action, insert.Timestamp}
// 	placeholder := []string{"$1", "$2", "$3", "$4"}
// 	stmt := "INSERT INTO device_history (" + strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholder, ", ") + `)
// 			 RETURNING id`
// 	if err := p.db.QueryRow(context.Background(), stmt, args...).Scan(
// 		&insert.ID,
// 	); err != nil {
// 		return nil, err
// 	}
// 	return insert, nil
// }

// func (p *_postgres) ListDeviceHistory(find *store.FindDeviceHistory) ([]*store.DeviceHistory, error) {
// 	where, args := []string{}, []any{}
// 	if v := find.DeviceID; v != nil {
// 		where, args = append(where, "device_id = $1"), append(args, *v)
// 	}
// 	if v := find.StartTime; v != nil {
// 		if len(args) > 0 {
// 			where = append(where, fmt.Sprintf("timestamp >= $%d", len(args)+1))
// 		} else {
// 			where = append(where, "timestamp >= $1")
// 		}
// 		args = append(args, *v)
// 	}
// 	if v := find.EndTime; v != nil {
// 		if len(args) > 0 {
// 			where = append(where, fmt.Sprintf("timestamp <= $%d", len(args)+1))
// 		} else {
// 			where = append(where, "timestamp <= $1")
// 		}
// 		args = append(args, *v)
// 	}
// 	stmt := `SELECT
// 				id,
// 				device_id,
// 				action,
// 				timestamp,
// 				user_id
// 			 FROM device_history`
// 	if len(where) > 0 {
// 		stmt += " WHERE " + strings.Join(where, " AND ")
// 	}
// 	stmt += " ORDER BY timestamp DESC"
// 	if v := find.Limit; v != nil {
// 		stmt += fmt.Sprintf(" LIMIT %d", *v)
// 	}
// 	rows, err := p.db.Query(context.Background(), stmt, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	list := make([]*store.DeviceHistory, 0)
// 	for rows.Next() {
// 		var history store.DeviceHistory
// 		if err := rows.Scan(
// 			&history.ID,
// 			&history.DeviceID,
// 			&history.Action,
// 			&history.Timestamp,
// 		); err != nil {
// 			return nil, err
// 		}
// 		list = append(list, &history)
// 	}
// 	return list, nil
// }

// func (p *_postgres) CreateDeviceSchedule(create *store.DeviceSchedule) (*store.DeviceSchedule, error) {
// 	fields := []string{"device_id", "action", "time", "repeat", "is_active"}
// 	args := []any{create.DeviceID, create.Action, create.Time, create.Repeat, create.IsActive}
// 	placeholder := []string{"$1", "$2", "$3", "$4", "$5"}
// 	stmt := "INSERT INTO device_schedules (" + strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholder, ", ") + `)
// 			 RETURNING id`
// 	if err := p.db.QueryRow(context.Background(), stmt, args...).Scan(
// 		&create.ID,
// 	); err != nil {
// 		return nil, err
// 	}
// 	return create, nil
// }
