package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/safe-homie/backend/internal/store"
)

func (p *_postgres) GetDeviceByID(find *store.FindDevice) (*store.Device, error) {
	where, args := []string{}, []any{}
	if v := find.ID; v != nil {
		where, args = append(where, "id = $1"), append(args, *v)
	}
	stmt := `SELECT
                id,
				serial_device,
                name,
                type,
				room
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
		&device.SerialDevice,
		&device.Name,
		&device.Type,
		&device.Room,
	); err != nil {
		return nil, err
	}
	return &device, nil
}

func (p *_postgres) ListDevices(find *store.FindDevice) ([]*store.Device, error) {
	where, args := []string{}, []any{}
	if v := find.Room; v != nil {
		where, args = append(where, "room = $1"), append(args, *v)
	}
	stmt := `SELECT
                id,
                serial_device,
                name,
                type,
				room
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
			&device.SerialDevice,
			&device.Name,
			&device.Type,
			&device.Room,
		); err != nil {
			return nil, err
		}
		list = append(list, &device)
	}
	return list, nil
}

func (p *_postgres) UpdateDevice(edit *store.UpdateDevice) (*store.UpdateDevice, error) {
	setClauses := []string{}
	args := []any{}
	argIdx := 1

	if edit.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *edit.Name)
		argIdx++
	}
	if edit.Room != nil {
		setClauses = append(setClauses, fmt.Sprintf("room = $%d", argIdx))
		args = append(args, *edit.Room)
		argIdx++
	}
	if edit.Type != nil {
		setClauses = append(setClauses, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, *edit.Type)
		argIdx++
	}

	if len(setClauses) == 0 {
		return edit, nil // không có gì để update
	}

	// Thêm ID để WHERE
	args = append(args, edit.ID)
	stmt := fmt.Sprintf("UPDATE devices SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIdx)

	_, err := p.db.Exec(context.Background(), stmt, args...)
	if err != nil {
		return nil, err
	}
	return edit, nil
}

func (p *_postgres) InsertStatusData(insert *store.DeviceStatus) error {
	stmt := `
	INSERT INTO device_status (device_id, active, schedule_enable, state, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err := p.db.Exec(context.Background(), stmt, insert.DeviceID, insert.Active, insert.ScheduleEnable, insert.State, insert.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert into device_status: %w", err)
	}
	return nil
}

func (p *_postgres) ListDeviceHistory(find *store.FindDeviceHistory) ([]*store.DeviceHistory, error) {
	where, args := []string{}, []any{}
	argIndex := 1

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

	stmt := `
		SELECT
			id,
			device_id,
			action,
			state,
			by,
			timestamp
		FROM device_history`

	if len(where) > 0 {
		stmt += " WHERE " + strings.Join(where, " AND ")
	}
	stmt += " ORDER BY timestamp DESC"

	if v := find.Limit; v != nil {
		stmt += fmt.Sprintf(" LIMIT %d", *v)
	}

	rows, err := p.db.Query(context.Background(), stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*store.DeviceHistory, 0)
	for rows.Next() {
		var history store.DeviceHistory
		var stateJSON []byte

		if err := rows.Scan(
			&history.ID,
			&history.DeviceID,
			&history.Action,
			&stateJSON,
			&history.By,
			&history.Timestamp,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(stateJSON, &history.State); err != nil {
			return nil, fmt.Errorf("failed to parse state json: %w", err)
		}

		list = append(list, &history)
	}

	return list, nil
}

func (p *_postgres) GetDeviceStatus(find *store.FindDevice) (*store.DeviceStatus, error) {
	if find.ID == nil {
		return nil, fmt.Errorf("device ID is required")
	}
	stmt := `
		SELECT
			device_id,
			active,
			schedule_enable,
			state,
			updated_at
		FROM device_status
		WHERE device_id = $1
	`

	var (
		deviceID       int32
		active         bool
		scheduleEnable bool
		stateData      []byte
		updatedAt      time.Time
	)

	err := p.db.QueryRow(context.Background(), stmt, *find.ID).Scan(
		&deviceID,
		&active,
		&scheduleEnable,
		&stateData,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy thiết bị: %v", err)
	}

	var state map[string]interface{}
	if err := json.Unmarshal(stateData, &state); err != nil {
		return nil, fmt.Errorf("lỗi giải mã state JSON: %v", err)
	}
	var schedule store.DeviceSchedule
	scheduleStmt := `
		SELECT
			id,
			device_id,
			action,
			scheduled_at,
			recurring,
			is_active,
			created_at,
			updated_at
		FROM device_schedule
		WHERE device_id = $1 AND is_active = true
		ORDER BY scheduled_at DESC
		LIMIT 1
	`
	err = p.db.QueryRow(context.Background(), scheduleStmt, *find.ID).Scan(
		&schedule.ID,
		&schedule.DeviceID,
		&schedule.Action,
		&schedule.ScheduledAt,
		&schedule.Recurring,
		&schedule.IsActive,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("lỗi lấy schedule: %v", err)
	}
	status := &store.DeviceStatus{
		DeviceID:       deviceID,
		Active:         active,
		ScheduleEnable: scheduleEnable,
		State:          state,
		Schedule:       schedule,
		UpdatedAt:      updatedAt,
	}

	return status, nil
}

func (p *_postgres) UpdateDeviceStatus(status *store.DeviceStatus) (*store.DeviceStatus, error) {
	stateJSON, err := json.Marshal(status.State)
	if err != nil {
		return nil, fmt.Errorf("lỗi mã hóa state: %v", err)
	}

	if status.UpdatedAt.IsZero() {
		status.UpdatedAt = time.Now()
	}

	stmt := `
			UPDATE device_status
			SET
				active = $1,
				schedule_enable = $2,
				state = $3,
				updated_at = $4
			WHERE device_id = $5
			RETURNING device_id
		`

	var id int32
	err = p.db.QueryRow(
		context.Background(),
		stmt,
		status.Active,
		status.ScheduleEnable,
		stateJSON,
		status.UpdatedAt,
		status.DeviceID,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("lỗi cập nhật trạng thái thiết bị: %v", err)
	}

	return status, nil
}

func (p *_postgres) CreateDeviceHistory(create *store.DeviceHistory) (*store.DeviceHistory, error) {
	stateJSON, err := json.Marshal(create.State)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal state: %w", err)
	}
	fields := []string{"device_id", "action", "state", "by", "timestamp"}
	args := []any{create.DeviceID, create.Action, stateJSON, create.By, create.Timestamp}
	placeholder := []string{"$1", "$2", "$3", "$4", "$5"}
	stmt := "INSERT INTO device_history (" + strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholder, ", ") + `)
	RETURNING id`

	if err := p.db.QueryRow(context.Background(), stmt, args...).Scan(&create.ID); err != nil {
		return nil, err
	}

	return create, nil
}

// func (p *_postgres) ListDeviceSchedule(find *store.FindDeviceSchedule) ([]*store.DeviceSchedule, error) {
// 	where, args := []string{}, []any{}
// 	argIndex := 1
// 	if v := find.DeviceID; v != nil {
// 		where = append(where, fmt.Sprintf("device_id = $%d", argIndex))
// 		args = append(args, *v)
// 		argIndex++
// 	}
// 	stmt := `SELECT
// 				id,
// 				device_id,
// 				action,
// 				CAST(start_date AS TIMESTAMPTZ) + time AS scheduled_at,
// 				repeat AS recurring,
// 				is_active,
// 				created_at,
// 				updated_at,
// 			 FROM device_schedules`
// 	if len(where) > 0 {
// 		stmt += " WHERE " + strings.Join(where, " AND ")
// 	}
// 	stmt += " ORDER BY scheduled_at ASC"
// 	if v := find.Limit; v != nil && *v > 0 {
// 		stmt += fmt.Sprintf(" LIMIT %d", *v)
// 	}

// 	rows, err := p.db.Query(context.Background(), stmt, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	list := make([]*store.DeviceSchedule, 0)
// 	for rows.Next() {
// 		var schedule store.DeviceSchedule
// 		if err := rows.Scan(
// 			&schedule.ID,
// 			&schedule.DeviceID,
// 			&schedule.Action,
// 			&schedule.ScheduledAt,
// 			&schedule.Recurring,
// 			&schedule.IsActive,
// 			&schedule.CreatedAt,
// 			&schedule.UpdatedAt,
// 		); err != nil {
// 			return nil, err
// 		}
// 		list = append(list, &schedule)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return list, nil
// }

// func (p *_postgres) CreateDeviceSchedule(create *store.DeviceSchedule) (*store.DeviceSchedule, error) {
// 	now := time.Now()
// 	create.CreatedAt = now
// 	create.UpdatedAt = now

// 	fields := []string{"device_id", "action", "scheduled_at", "recurring", "is_active", "created_at", "updated_at"}
// 	args := []any{create.DeviceID, create.Action, create.ScheduledAt, create.Recurring, create.IsActive, create.CreatedAt, create.UpdatedAt}
// 	placeholder := []string{"$1", "$2", "$3", "$4", "$5", "$6", "$7"}

// 	stmt := "INSERT INTO device_schedules (" + strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholder, ", ") + `)
//              RETURNING id`

// 	if err := p.db.QueryRow(context.Background(), stmt, args...).Scan(&create.ID); err != nil {
// 		return nil, err
// 	}

// 	return create, nil
// }

// func (p *_postgres) UpdateDeviceSchedule(edit *store.UpdateDeviceSchedule) (*store.DeviceSchedule, error) {
// 	var setClauses []string
// 	var args []any
// 	argPos := 1

// 	if edit.Action != nil {
// 		setClauses = append(setClauses, fmt.Sprintf("action = $%d", argPos))
// 		args = append(args, *edit.Action)
// 		argPos++
// 	}
// 	if edit.ScheduledAt != nil {
// 		setClauses = append(setClauses, fmt.Sprintf("scheduled_at = $%d", argPos))
// 		args = append(args, *edit.ScheduledAt)
// 		argPos++
// 	}
// 	if edit.Recurring != nil {
// 		setClauses = append(setClauses, fmt.Sprintf("recurring = $%d", argPos))
// 		args = append(args, *edit.Recurring)
// 		argPos++
// 	}

// 	now := time.Now()
// 	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argPos))
// 	args = append(args, now)
// 	edit.UpdatedAt = &now
// 	argPos++

// 	if len(setClauses) == 0 {
// 		return nil, fmt.Errorf("no fields to update")
// 	}

// 	args = append(args, edit.ID)
// 	stmt := `
//         UPDATE device_schedules
//         SET ` + strings.Join(setClauses, ", ") + `
//         WHERE id = $` + fmt.Sprint(argPos) + `
//         RETURNING id, device_id, action, scheduled_at, recurring, is_active, created_at, updated_at
//     `

// 	var updated store.DeviceSchedule
// 	if err := p.db.QueryRow(context.Background(), stmt, args...).Scan(
// 		&updated.ID,
// 		&updated.DeviceID,
// 		&updated.Action,
// 		&updated.ScheduledAt,
// 		&updated.Recurring,
// 		&updated.IsActive,
// 		&updated.CreatedAt,
// 		&updated.UpdatedAt,
// 	); err != nil {
// 		return nil, err
// 	}

// 	return &updated, nil
// }

// func (p *_postgres) DeleteDeviceSchedule(find *store.FindDevice) error {
// 	return nil
// }
