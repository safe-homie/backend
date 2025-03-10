package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/safe-homie/backend/internal/store"
)

func (p *_postgres) CreateSensor(create *store.Sensor) (*store.Sensor, error) {
	fields := []string{"type", "location", "unit", "name", "threshold_warning", "threshold_danger"}
	args := []any{create.Type, create.Location, create.Unit, create.Name, create.ThresholdWarning, create.ThresholdDanger}
	placeholder := []string{"$1", "$2", "$3", "$4", "$5", "$6"}
	stmt := "INSERT INTO sensors (" + strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholder, ", ") + `) 
			 RETURNING id`
	if err := p.db.QueryRow(context.Background(), stmt, args...).Scan(
		&create.ID,
	); err != nil {
		return nil, err
	}
	return create, nil
}

func (p *_postgres) UpdateSensor(update *store.UpdateSensor) (*store.Sensor, error) {
	set, args := []string{}, []any{}
	if v := update.Type; v != nil {
		set, args = append(set, "type = $1"), append(args, *v)
	}
	if v := update.Location; v != nil {
		set, args = append(set, "location = $2"), append(args, *v)
	}
	if v := update.Unit; v != nil {
		set, args = append(set, "type = $3"), append(args, *v)
	}
	if v := update.ThresholdWarning; v != nil {
		set, args = append(set, "threshold_warning = $4"), append(args, *v)
	}
	if v := update.ThresholdDanger; v != nil {
		set, args = append(set, "threshold_danger = $5"), append(args, *v)
	}
	if v := update.Name; v != nil {
		set, args = append(set, "name = $6"), append(args, *v)
	}
	args = append(args, update.ID)
	stmt := `UPDATE sensors
			 SET` + strings.Join(set, ", ") + `
			 WHERE id = $7
			 RETURNING id, type, location, unit, name, threshold_warning, threshold_danger
			 `
	sensor := &store.Sensor{}
	if err := p.db.QueryRow(context.Background(), stmt, args...).Scan(
		&sensor.ID,
		&sensor.Type,
		&sensor.Location,
		&sensor.Unit,
		&sensor.Name,
		&sensor.ThresholdWarning,
		&sensor.ThresholdDanger,
	); err != nil {
		return nil, err
	}
	return sensor, nil
}

func (p *_postgres) ListSensors(find *store.FindSensor) ([]*store.Sensor, error) {
	where, args := []string{}, []any{}
	if v := find.Location; v != nil {
		where, args = append(where, "location = $1"), append(args, *v)
	}
	stmt := `SELECT 
					id,
					type,
					location,
					unit,
					name,
					threshold_warning,
					threshold_danger
			 FROM sensors`
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
	list := make([]*store.Sensor, 0)
	for rows.Next() {
		var sensor store.Sensor
		if err := rows.Scan(
			&sensor.ID,
			&sensor.Type,
			&sensor.Location,
			&sensor.Unit,
			&sensor.Name,
			&sensor.ThresholdWarning,
			&sensor.ThresholdDanger,
		); err != nil {
			return nil, err
		}
		list = append(list, &sensor)
	}
	return list, nil
}

func (p *_postgres) GetSensorByID(find *store.FindSensor) (*store.Sensor, error) {
	where, args := []string{}, []any{}
	if v := find.ID; v != nil {
		where, args = append(where, "id = $1"), append(args, *v)
	}
	stmt := `SELECT 
					id,
					type,
					location,
					unit,
					name,
					threshold_warning,
					threshold_danger
			 FROM sensors`
	if len(where) > 0 {
		stmt += " WHERE " + strings.Join(where, " AND ")
	}
	if v := find.Limit; v != nil {
		stmt += fmt.Sprintf(" LIMIT %d", *v)
	}
	var sensor store.Sensor
	if err := p.db.QueryRow(context.Background(), stmt, args...).Scan(
		&sensor.ID,
		&sensor.Type,
		&sensor.Location,
		&sensor.Unit,
		&sensor.Name,
		&sensor.ThresholdWarning,
		&sensor.ThresholdDanger,
	); err != nil {
		return nil, err
	}
	return &sensor, nil
}

func (p *_postgres) InsertSensorData(insert *store.SensorData) (*store.SensorData, error) {
	fields := []string{"sensor_id", "time", "value"}
	args := []any{insert.SensorID, insert.Time, insert.Value}
	placeholder := []string{"$1", "$2", "$3"}
	stmt := "INSERT INTO sensor_data (" + strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholder, ", ") + `)`
	if err := p.db.QueryRow(context.Background(), stmt, args...).Scan(); err != nil {
		return nil, err
	}
	return insert, nil
}

func (p *_postgres) GetLatestSensorData(find *store.FindSensorData) (*store.SensorData, error) {
	stmt := `SELECT sensor_id, value, time 
			 FROM sensor_data
			 WHERE sensor_id = $1
			 ORDER BY time DESC
			 LIMIT 1
			 `
	args := []any{find.SensorID}
	var data store.SensorData
	row := p.db.QueryRow(context.Background(), stmt, args...)
	if err := row.Scan(
		&data.SensorID,
		&data.Value,
		&data.Time,
	); err != nil {
		return nil, err
	}
	return &data, nil
}

func (p *_postgres) ListLatestSensorDataByLocation(find *store.FindSensorData) ([]*store.SensorDataWithProfile, error) {
	stmt := `
		SELECT 
			s.id AS sensor_id, 
			s.type, 
			s.name, 
			s.location, 
			s.unit, 
			sd.value, 
			sd.time
		FROM sensors s
		LEFT JOIN LATERAL (
			SELECT sd.value, sd.time
			FROM sensor_data sd
			WHERE sd.sensor_id = s.id
			ORDER BY sd.time DESC
			LIMIT 1
		) sd ON true
		WHERE s.location = $1;
	`
	rows, err := p.db.Query(context.Background(), stmt, find.Location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*store.SensorDataWithProfile
	for rows.Next() {
		var data store.SensorDataWithProfile
		if err := rows.Scan(
			&data.SensorID,
			&data.Type,
			&data.Name,
			&data.Location,
			&data.Unit,
			&data.Value,
			&data.Time,
		); err != nil {
			return nil, err
		}
		list = append(list, &data)
	}
	return list, nil
}

func (p *_postgres) GetAverageSensorData() (*store.SensorData, error) {
	return &store.SensorData{}, nil
}
