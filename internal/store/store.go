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
	ListDeviceHistory(find *FindDeviceHistory) ([]*DeviceHistory, error)
	ListDeviceSchedule(find *FindDeviceSchedule) ([]*DeviceSchedule, error)

	CreateDevice(create *Device) (*Device, error)
	UpdateDevice(update *UpdateDevice) (*UpdateDevice, error)

	// CreateDeviceSchedule(create *DeviceSchedule) (*DeviceSchedule, error)
	// InsertDeviceHistory(insert *DeviceHistory) (*DeviceHistory, error)
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
