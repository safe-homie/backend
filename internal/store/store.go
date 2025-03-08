package store

type store struct {
	driver Driver
}

type Store interface {
	Migrate() error
	Close()

	CreateSensor(create *Sensor) (*Sensor, error)
	UpdateSensor(update *UpdateSensor) (*Sensor, error)
	InsertSensorData(insert *SensorData) (*SensorData, error)
	GetLatestSensorData(find *FindSensorData) (*SensorData, error)
}

func New(driver Driver) *store {
	return &store{driver: driver}
}

func (s *store) Migrate() error {
	return s.driver.Migrate()
}

func (s *store) Close() {
	s.driver.Close()
}
