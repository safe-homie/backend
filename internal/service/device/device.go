package device

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/safe-homie/backend/internal/cache"
	"github.com/safe-homie/backend/internal/domain"
	"github.com/safe-homie/backend/internal/store"
	"github.com/safe-homie/backend/internal/util"
)

// TODO: Implement Device service
// Including REST + MQTT
type DeviceService interface {
	GetDevice(id int32) (*store.Device, error)
	ListDevices(location string) ([]*store.Device, error)
	ListDeviceHistory(id int32, req *domain.GetDeviceHistoryRequest) ([]*store.DeviceHistory, error)
	ListDeviceSchedule(deviceID int32) ([]*store.DeviceSchedule, error)

	// ControlDevice(id int32, status string) (*store.Device, error)
	CreateDevice(req *domain.CreateDeviceRequest) (*store.Device, error)
	UpdateDevice(id int32, req *domain.UpdateDeviceRequest) (*store.UpdateDevice, error) // will change later
	// CreateDeviceSchedule(req *domain.CreateScheduleRequest) (*store.DeviceSchedule, error)

	// ApplyPreset()
	// SetAutoMode()
}

func NewService(store store.Store) DeviceService {
	srv := &deviceService{store: store, cache: cache.NewInMemoryCache()}
	srv.loadCacheFromDB()
	return srv
}

type deviceService struct {
	store store.Store
	cache cache.Cache
}

func (s *deviceService) loadCacheFromDB() {
	devices, err := s.store.ListDevices(&store.FindDevice{})
	if err != nil {
		return
	}
	for _, device := range devices {
		s.cache.Set(strconv.Itoa(int(device.ID)), device.Type, 0)
	}
}

func (s *deviceService) GetDevice(id int32) (*store.Device, error) {
	find := store.FindDevice{ID: &id}
	deviceDB, err := s.store.GetDevice(&find)
	if err != nil {
		return nil, err
	}
	return deviceDB, nil
}

func (s *deviceService) ListDevices(location string) ([]*store.Device, error) {
	location = util.GetValueOrDefault(location, string(domain.DefaultDeviceLocation))
	find := store.FindDevice{Location: &location}
	devices, err := s.store.ListDevices(&find)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (s *deviceService) CreateDevice(req *domain.CreateDeviceRequest) (*store.Device, error) {
	create := store.Device{
		Name:     req.Name,
		Location: req.Location,
		Type:     req.Type,
		Power:    req.Power,
		Level:    req.Level,
	}
	createDB, err := s.store.CreateDevice(&create)
	if err != nil {
		return nil, err
	}
	return createDB, nil
}

func (s *deviceService) UpdateDevice(id int32, req *domain.UpdateDeviceRequest) (*store.UpdateDevice, error) {
	update := store.UpdateDevice{
		Location: &req.Location,
		Name:     &req.Name,
		Type:     &req.Type,
	}
	updateDB, err := s.store.UpdateDevice(&update)
	if err != nil {
		return nil, err
	}
	return updateDB, nil
}

func (s *deviceService) ListDeviceHistory(id int32, req *domain.GetDeviceHistoryRequest) ([]*store.DeviceHistory, error) {
	// mặc định là 24h
	now := time.Now()
	startTime := util.GetValueOrDefault(req.StartTime, now.Add(-24*time.Hour).Format(time.RFC3339))
	endTime := util.GetValueOrDefault(req.EndTime, now.Format(time.RFC3339))
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, fmt.Errorf("lỗi định dạng thời gian bắt đầu: %v", err)
	}
	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return nil, fmt.Errorf("lỗi định dạng thời gian kết thúc: %v", err)
	}
	if start.After(end) {
		return nil, errors.New("khoảng thời gian không hợp lệ: thời gian bắt đầu muộn hơn thời gian kết thúc")
	}
	find := store.FindDeviceHistory{
		DeviceID:  &id,
		StartTime: &start,
		EndTime:   &end,
	}
	histories, err := s.store.ListDeviceHistory(&find)
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (s *deviceService) ListDeviceSchedule(idDevice int32) ([]*store.DeviceSchedule, error) {
	find := store.FindDeviceSchedule{DeviceID: &idDevice}
	schedules, err := s.store.ListDeviceSchedule(&find)
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

// func (s *deviceService) UpdateDevice(id int32, req *domain.UpdateDeviceRequest) error {
// 	existingDevice, err := s.store.GetDeviceByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if existingDevice == nil {
// 		return nil, fmt.Errorf("device not found")
// 	}

// 	updateData := store.Device{
// 		ID:       id,
// 		Name:     *req.Name,
// 		Location: *req.Location,
// 		Type:     *req.Type,
// 		Power:    *req.Power,
// 		Level:    *req.Level,
// 	}

// }

// func (s *deviceService) UpdateDevice(id int32, req *domain.UpdateDeviceRequest) (*store.Device, error) {

// 	update := domain.UpdateDeviceRequest{
// 		ID: id,
// 	}

// 	if req.Name != nil {
// 		update.Name = req.Name
// 	}
// 	if req.Location != nil {
// 		update.Location = req.Location
// 	}
// 	if req.Type != nil {
// 		update.Type = req.Type
// 	}
// 	if req.Power != nil {
// 		update.Power = req.Power
// 	}
// 	if req.Level != nil {

// 		levelStr := fmt.Sprintf("%d", *req.Level)
// 		update.Level = &levelStr
// 	}

// 	updatedDevice, err := s.store.UpdateDevice(&update)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update device in store: %w", err)
// 	}

// 	return updatedDevice, nil
// }

// func (s *deviceService) ControlDevice(id int32, status string) (*store.Device, error) {
// 	return nil, nil
// }

// func (s *deviceService) GetDeviceHistory(id int32, startTime string, endTime string) ([]*store.DeviceHistory, error) {
// 	return nil, nil
// }
// func (s *deviceService) CreateDeviceSchedule(req *domain.CreateScheduleRequest) (*store.DeviceSchedule, error) {
// 	return nil, nil
// }

// func (s *deviceService) ListDeviceSchedules(deviceID int32) ([]*store.DeviceSchedule, error) {
// 	return nil, nil
// }

// // func (s *deviceService) GetDeviceHistory(id int32, startTime string, endTime string) ([]*store.DeviceHistory, error) {
// // 	// find:=store
// // 	// if startTime != "" {
// // 	// 	start, _ := time.Parse(time.RFC3339, startTime)
// // 	// 	find.StartTime = &start
// // 	// }
// // 	// if endTime != "" {
// // 	// 	end, _ := time.Parse(time.RFC3339, endTime)
// // 	// 	find.EndTime = &end
// // 	// }
// // 	// history, err := s.store.ListDeviceHistory(&find)
// // 	// if err != nil {
// // 	// 	return nil, err
// // 	// }
// // 	// return domain.ToDeviceHistoryListDTO(history), nil
// // 	return nil, nil
// // }

// func (s *deviceService) getDeviceType(deviceID int32) (string, error) {
// 	if val, ok := s.cache.Get(strconv.Itoa(int(deviceID))); ok {
// 		return val.(string), nil
// 	}
// 	device, err := s.store.GetDevice(&store.FindDevice{ID: &deviceID})
// 	if err != nil {
// 		return "", err
// 	}
// 	s.cache.Set(strconv.Itoa(int(device.ID)), device.Type, 0)
// 	return device.Type, nil
// }

// func (s *deviceService) validateDeviceData(data *domain.DeviceDataMessage) bool {
// 	sensorType, _ := s.getSensorType(data.ID)
// 	return sensorType == data.Type
// }
