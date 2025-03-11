package domain

const (
	SensorHumidity    = "humidity"
	SensorTemperature = "temperature"
	SensorLight       = "light"

	DefaultThresholdWarning = 50.0
	DefaultThresholdDanger  = 80.0

	UnitHumidity    = "%"
	UnitTemperature = "°C"
	UnitLight       = "lx"
	UnitUnknown     = "unknown"

	DefaultSensorsLocation = "living room"
)

func IsValidSensorType(sensorType string) bool {
	switch sensorType {
	case SensorHumidity, SensorTemperature, SensorLight:
		return true
	default:
		return false
	}
}

func GetSensorUnit(sensorType string) string {
	switch sensorType {
	case SensorHumidity:
		return UnitHumidity
	case SensorLight:
		return UnitLight
	case SensorTemperature:
		return UnitTemperature
	default:
		return UnitUnknown
	}
}
