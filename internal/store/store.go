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
}

type StatusHistoryOperations interface {
	// GetDeviceStatus(find *FindDevice) (*DeviceStatus, error)
	// UpdateDeviceStatus(status *DeviceStatus) (*DeviceStatus, error)

	// CreateDeviceHistory(create *DeviceHistory) (*DeviceHistory, error)

	GetStatus(deviceID int32) (*DeviceStatus, error)
	// UpdateStatus(status *DeviceStatus) error
	// RecordHistory(history *DeviceHistory) error
}

// type ScheduleOperations interface {
// 	ListDeviceSchedule(find *FindDeviceSchedule) ([]*DeviceSchedule, error)
// 	CreateDeviceSchedule(create *DeviceSchedule) (*DeviceSchedule, error)
// 	UpdateDeviceSchedule(edit *UpdateDeviceSchedule) (*DeviceSchedule, error)
// 	DeleteDeviceSchedule(find *FindDevice) error
// }

func New(driver Driver) Store {
	return &store{driver: driver}
}

func (s *store) Migrate() error {
	return s.driver.Migrate()
}

func (s *store) Close() {
	s.driver.Close()
}
