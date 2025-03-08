package postgres

import (
	"context"
	"strings"

	"github.com/safe-homie/backend/internal/store"
)

func (p *_postgres) CreateSensor(create *store.Sensor) (*store.Sensor, error) {
	fields := []string{"type", "location", "unit", "threshold_warning", "threshold_danger"}
	args := []any{create.Type, create.Location, create.Unit, create.ThresholdWarning, create.ThresholdDanger}
	placeholder := []string{"$1", "$2", "$3", "$4", "$5"}
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
	args = append(args, update.ID)
	stmt := `UPDATE sensors
			 SET` + strings.Join(set, ", ") + `
			 WHERE id = $6
			 RETURNING id, type, location, unit, threshold_warning, threshold_danger
			 `
	sensor := &store.Sensor{}
	if err := p.db.QueryRow(context.Background(), stmt, args...).Scan(
		&sensor.ID,
		&sensor.Type,
		&sensor.Location,
		&sensor.Unit,
		&sensor.ThresholdWarning,
		&sensor.ThresholdDanger,
	); err != nil {
		return nil, err
	}
	return sensor, nil
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

func (p *_postgres) GetAverageSensorData() (*store.SensorData, error) {
	return &store.SensorData{}, nil
}
