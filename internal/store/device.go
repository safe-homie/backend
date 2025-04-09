package store

import "time"

type (
	Device struct {
		ID       int32
		Location string
		Name     string
		Type     string
		Power    string
		Level    string
		CreateAt time.Time
		UpdateAt time.Time
	}
	FindDevice struct {
		ID       *int32
		Location *string
		Limit    *int
	}
	UpdateDevice struct {
		ID        int32
		Location  *string
		Name      *string
		Type      *string
		UpdatedAt *time.Time
	}
	DeviceHistory struct {
		ID        int32
		DeviceID  int32
		Action    string
		Timestamp time.Time
	}
	DeviceSchedule struct {
		ID          int32
		DeviceID    int32
		Action      string // "TURN_ON", "TURN_OFF"
		ScheduledAt time.Time
		Recurring   string // "DAILY", "WEEKLY", "ONCE"
		IsActive    bool
		CreatedAt   time.Time
	}
	FindDeviceHistory struct {
		DeviceID  *int32
		StartTime *time.Time
		EndTime   *time.Time
		Limit     *int
	}
	FindDeviceSchedule struct {
		DeviceID *int32
		Limit    *int
	}
)

func (s *store) ListDevices(find *FindDevice) ([]*Device, error) {
	return s.driver.ListDevices(find)
}
func (s *store) GetDevice(find *FindDevice) (*Device, error) {
	return s.driver.GetDeviceByID(find)
}

func (s *store) ListDeviceHistory(find *FindDeviceHistory) ([]*DeviceHistory, error) {
	return s.driver.ListDeviceHistory(find)
}
func (s *store) ListDeviceSchedule(find *FindDeviceSchedule) ([]*DeviceSchedule, error) {
	return s.driver.ListDeviceSchedule(find)
}

func (s *store) CreateDevice(create *Device) (*Device, error) {
	return s.driver.CreateDevice(create)
}
func (s *store) UpdateDevice(update *UpdateDevice) (*UpdateDevice, error) {
	return s.driver.UpdateDevice(update)
}

// func (s *store) InsertDeviceHistory(insert *DeviceHistory) (*DeviceHistory, error) {
// 	return s.driver.InsertDeviceHistory(insert)
// }

// func (s *store) CreateDeviceSchedule(create *DeviceSchedule) (*DeviceSchedule, error) {
// 	return s.driver.CreateDeviceSchedule(create)
// }
