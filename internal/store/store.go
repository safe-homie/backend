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

	GetDevice(find *FindDevice) (*Device, error)
	ListDevices(find *FindDevice) ([]*Device, error)
	UpdateDevice(edit *UpdateDevice) (*UpdateDevice, error)
	ListDeviceHistory(find *FindDeviceHistory) ([]*DeviceHistory, error)
	InsertStatusData(insert *DeviceStatus) error
	GetStatus(find *FindDevice) (*DeviceStatus, error)
	UpdateStatus(status *DeviceStatus) (*DeviceStatus, error)
	RecordHistory(history *DeviceHistory) (*DeviceHistory, error)

	SaveToken(save *Token) (*Token, error)
	GetToken(find *FindToken) (*Token, error)
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
