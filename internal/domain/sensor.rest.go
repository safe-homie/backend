package domain

import (
	"time"

	"github.com/safe-homie/backend/internal/store"
)

type GetSensorResponse struct {
	ID               int32   `json:"id"`
	Type             string  `json:"type"`
	Location         string  `json:"location"`
	Name             string  `json:"name"`
	Unit             string  `json:"unit"`
	ThresholdWarning float64 `json:"threshold_warning"`
	ThresholdDanger  float64 `json:"threshold_danger"`
}

type CreateSensorRequest struct {
	Type             string   `json:"type" validate:"required,sensor_type"`
	Location         string   `json:"location" validate:"required"`
	Name             string   `json:"name" validate:"required"`
	ThresholdWarning *float64 `json:"threshold_warning,omitempty" validate:"omitempty,gte=0"`
	ThresholdDanger  *float64 `json:"threshold_danger,omitempty" validate:"omitempty,gte=0"`
}

type GetSensorDataResponse struct {
	SensorID int32     `json:"id"`
	Value    float64   `json:"value" validate:"gt=0"`
	Time     time.Time `json:"time"`
}

type GetSensorDataWithProfileResponse struct {
	SensorID int32     `json:"id"`
	Type     string    `json:"type"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Unit     string    `json:"unit"`
	Value    float64   `json:"value"`
	Time     time.Time `json:"time"`
}

func ToSensorDTO(s *store.Sensor) *GetSensorResponse {
	return &GetSensorResponse{
		ID:               s.ID,
		Type:             s.Type,
		Location:         s.Location,
		Unit:             s.Unit,
		Name:             s.Name,
		ThresholdWarning: s.ThresholdWarning,
		ThresholdDanger:  s.ThresholdDanger,
	}
}

func ToSensorsDTO(in []*store.Sensor) []*GetSensorResponse {
	out := make([]*GetSensorResponse, 0)
	for _, item := range in {
		sensor := ToSensorDTO(item)
		out = append(out, sensor)
	}
	return out
}

func ToSensorDataDTO(s *store.SensorData) *GetSensorDataResponse {
	return &GetSensorDataResponse{
		SensorID: s.SensorID,
		Value:    s.Value,
		Time:     s.Time,
	}
}

func ToSensorDataListDTO(in []*store.SensorDataWithProfile) []*GetSensorDataWithProfileResponse {
	out := make([]*GetSensorDataWithProfileResponse, 0)
	for _, item := range in {
		data := GetSensorDataWithProfileResponse{
			SensorID: item.SensorID,
			Type:     item.Type,
			Name:     item.Name,
			Unit:     item.Unit,
			Location: item.Location,
			Value:    item.Value,
			Time:     item.Time,
		}
		out = append(out, &data)
	}
	return out
}
