package domain

import "time"

type (
	SensorDataMessage struct {
		ID    int32     `json:"id"`
		Value float64   `json:"value"`
		Time  time.Time `json:"time"`
	}
)
