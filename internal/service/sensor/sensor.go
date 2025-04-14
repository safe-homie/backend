package sensor

import (
	"fmt"
	"strconv"

	"github.com/safe-homie/backend/internal/cache"
	"github.com/safe-homie/backend/internal/domain"
	"github.com/safe-homie/backend/internal/store"
	"github.com/safe-homie/backend/internal/util"
	"github.com/safe-homie/backend/pkg/event"
)

// TODO: Implement Sensor service
// Including REST + MQTT
type SensorService interface {
	GetSensor(id int32) (*store.Sensor, error)
	ListSensors(location string) ([]*store.Sensor, error)
	CreateSensor(req *domain.CreateSensorRequest) (*store.Sensor, error)
	GetLatestSensorData(id int32) (*store.SensorData, error) // Not implemented
	InsertSensorData(data *domain.SensorDataMessage) (*store.SensorData, error)
	ListLatestSensorDataByLocation(location string) ([]*store.SensorDataWithProfile, error)
}

func NewService(store store.Store, evm event.EventManager) SensorService {
	srv := &sensorService{
		store:        store,
		cache:        cache.NewInMemoryCache(),
		eventManager: evm,
	}
	srv.loadCacheFromDB()
	return srv
}

type sensorService struct {
	store        store.Store
	cache        cache.Cache
	eventManager event.EventManager
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
	location = util.GetValueOrDefault(location, string(domain.DefaultSensorsLocation))
	find := store.FindSensor{Location: &location}
	sensors, err := s.store.ListSensors(&find)
	if err != nil {
		return nil, err
	}
	return sensors, nil
}

func (s *sensorService) CreateSensor(req *domain.CreateSensorRequest) (*store.Sensor, error) {
	create := store.Sensor{
		Type:     req.Type,
		Location: req.Location,
		Name:     req.Name,
	}
	create.Unit = domain.GetSensorUnit(req.Type)
	create.ThresholdWarning = util.GetPointerValueOrDefault(req.ThresholdWarning, domain.DefaultThresholdWarning)
	create.ThresholdDanger = util.GetPointerValueOrDefault(req.ThresholdDanger, domain.DefaultThresholdDanger)
	createDB, err := s.store.CreateSensor(&create)
	if err != nil {
		return nil, err
	}
	return createDB, nil
}

func (s *sensorService) ListLatestSensorDataByLocation(location string) ([]*store.SensorDataWithProfile, error) {
	location = util.GetValueOrDefault(location, string(domain.DefaultSensorsLocation))
	data, err := s.store.ListLatestSensorDataByLocation(&store.FindSensorData{Location: location})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *sensorService) InsertSensorData(data *domain.SensorDataMessage) (*store.SensorData, error) {
	valid := s.validateSensorData(data)
	if !valid {
		return nil, fmt.Errorf("sensor type and sensor id mismatch")
	}
	insert := store.SensorData{
		SensorID: data.ID,
		Value:    data.Value,
		Time:     data.Time,
	}
	insertDB, err := s.store.InsertSensorData(&insert)
	if err != nil {
		return nil, err
	}
	go s.handleIfExceeded(&insert)
	return insertDB, nil
}

// This method is not necessary at this time
func (s *sensorService) GetLatestSensorData(id int32) (*store.SensorData, error) {
	return &store.SensorData{}, nil
}

func (s *sensorService) loadCacheFromDB() {
	sensors, err := s.store.ListSensors(&store.FindSensor{})
	if err != nil {
		return
	}
	for _, sensor := range sensors {
		s.cache.Set(strconv.Itoa(int(sensor.ID)), sensor.Type, 0)
	}
}

func (s *sensorService) getSensorType(sensorID int32) (string, error) {
	if val, ok := s.cache.Get(strconv.Itoa(int(sensorID))); ok {
		return val.(string), nil
	}
	sensor, err := s.store.GetSensor(&store.FindSensor{ID: &sensorID})
	if err != nil {
		return "", err
	}
	s.cache.Set(strconv.Itoa(int(sensor.ID)), sensor.Type, 0)
	return sensor.Type, nil
}

func (s *sensorService) validateSensorData(data *domain.SensorDataMessage) bool {
	sensorType, _ := s.getSensorType(data.ID)
	return sensorType == data.Type
}

func (s *sensorService) handleIfExceeded(data *store.SensorData) {
	sensorDB, _ := s.store.GetSensor(&store.FindSensor{ID: &data.SensorID})
	if s.isExceeded(sensorDB.ThresholdDanger, data.Value) {
		ev := domain.SensorThresholdExceedEvent{
			Notify: domain.Notification{
				Title: "Cảnh báo vượt ngưỡng",
				Body: fmt.Sprintf("Cảm biến %s tại %s đo được giá trị %.2f, vượt ngưỡng an toàn",
					sensorDB.Name, sensorDB.Location, data.Value),
			},
		}
		s.eventManager.EmitEvent(domain.SensorThresholdExceed, ev)
	}
}

func (s *sensorService) isExceeded(threshold, value float64) bool {
	return value >= threshold
}
