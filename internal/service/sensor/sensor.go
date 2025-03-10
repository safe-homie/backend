package sensor

import (
	"github.com/safe-homie/backend/internal/domain"
	"github.com/safe-homie/backend/internal/store"
	"github.com/safe-homie/backend/internal/transport/rest/dto"
)

// TODO: Implement Sensor service
// Including REST + MQTT
type SensorService interface {
	GetSensor(id int32) (*store.Sensor, error)
	ListSensors(location string) ([]*store.Sensor, error)
	CreateSensor(req *dto.CreateSensorRequest) (*store.Sensor, error)
	GetLatestSensorData(id int32) (*store.SensorData, error)
	ListLatestSensorDataByLocation(location string) ([]*store.SensorDataWithProfile, error)
}

func NewService(store store.Store) SensorService {
	return &sensorService{store: store}
}

type sensorService struct {
	store store.Store
}

func (s *sensorService) GetSensor(id int32) (*store.Sensor, error) {
	find := store.FindSensor{ID: &id}
	sensorDB, err := s.store.GetSensor(&find)
	if err != nil {
		return nil, err
	}
	return sensorDB, nil
}

func (s *sensorService) ListSensors(location string) ([]*store.Sensor, error) {
	if location == "" {
		location = domain.DefaultSensorsLocation
	}
	find := store.FindSensor{Location: &location}
	sensors, err := s.store.ListSensors(&find)
	if err != nil {
		return nil, err
	}
	return sensors, nil
}

func (s *sensorService) CreateSensor(req *dto.CreateSensorRequest) (*store.Sensor, error) {
	create := store.Sensor{
		Type:     req.Type,
		Location: req.Location,
		Name:     req.Name,
	}
	create.Unit = domain.GetSensorUnit(req.Type)
	create.ThresholdWarning = getOrDefaultThresholdWarning(req.ThresholdWarning)
	create.ThresholdDanger = getOrDefaultThresholdDanger(req.ThresholdDanger)
	createDB, err := s.store.CreateSensor(&create)
	if err != nil {
		return nil, err
	}
	return createDB, nil
}

func (s *sensorService) ListLatestSensorDataByLocation(location string) ([]*store.SensorDataWithProfile, error) {
	if location == "" {
		location = domain.DefaultSensorsLocation
	}
	data, err := s.store.ListLatestSensorDataByLocation(&store.FindSensorData{Location: location})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// This method is not necessary at this time
func (s *sensorService) GetLatestSensorData(id int32) (*store.SensorData, error) {
	return &store.SensorData{}, nil
}

func getOrDefaultThresholdWarning(value *float64) float64 {
	if value != nil {
		return *value
	}
	return domain.DefaultThresholdWarning
}

func getOrDefaultThresholdDanger(value *float64) float64 {
	if value != nil {
		return *value
	}
	return domain.DefaultThresholdDanger
}
