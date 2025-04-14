package store

type store struct {
	driver Driver
}

type Store interface {
	Migrate() error
	Close()

	GetSensor(find *FindSensor) (*Sensor, error)
	ListSensors(find *FindSensor) ([]*Sensor, error)
	CreateSensor(create *Sensor) (*Sensor, error)
	UpdateSensor(update *UpdateSensor) (*Sensor, error)
	InsertSensorData(insert *SensorData) (*SensorData, error)
	GetLatestSensorData(find *FindSensorData) (*SensorData, error)
	ListLatestSensorDataByLocation(find *FindSensorData) ([]*SensorDataWithProfile, error)
	SaveToken(save *Token) (*Token, error)
}

func New(driver Driver) Store {
	return &store{driver: driver}
}

func (s *store) Migrate() error {
	return s.driver.Migrate()
}

func (s *store) Close() {
	s.driver.Close()
}
