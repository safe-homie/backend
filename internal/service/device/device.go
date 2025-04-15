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

type DeviceRepository interface {
	GetStatus(deviceID int32) (*store.DeviceStatus, error)
	UpdateStatus(status *store.DeviceStatus) error
	RecordHistory(history *store.DeviceHistory) error
}
type DeviceCommand interface {
	Execute() error
	GetDeviceID() int32
	GetCommandType() string
	GetState() map[string]interface{}
}

// TODO: Implement Device service
// Including REST + MQTT
type DeviceService interface {
	GetDevice(id int32) (*store.Device, error)
	ListDevices(room string) ([]*store.Device, error)
	UpdateDevice(id int32, req *domain.UpdateDeviceRequest) (*store.UpdateDevice, error)

	GetDeviceStatus(deviceID int32) (*store.DeviceStatus, error)

	ListDeviceHistory(id int32, req *domain.GetDeviceHistoryRequest) ([]*store.DeviceHistory, error)

	ExecuteCommand(cmd domain.DeviceCommand) error

	InsertStatusData(data *domain.DeviceMessage) error
	HandleControlMessage(data *domain.DeviceMessage) error

	// CreateDeviceHistory(id int32) (*store.DeviceHistory, error)
	// ListDeviceSchedule(deviceID int32) ([]*store.DeviceSchedule, error)
	// CreateDevice(req *domain.CreateDeviceRequest) (*store.Device, error)
	// CreateDeviceSchedule(req *domain.CreateScheduleRequest) (*store.DeviceSchedule, error)

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

func (s *deviceService) ListDevices(room string) ([]*store.Device, error) {
	room = util.GetValueOrDefault(room, string(domain.DefaultDeviceLocation))
	find := store.FindDevice{Room: &room}
	devices, err := s.store.ListDevices(&find)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (s *deviceService) UpdateDevice(id int32, req *domain.UpdateDeviceRequest) (*store.UpdateDevice, error) {
	update := store.UpdateDevice{
		ID:   req.ID,
		Room: &req.Room,
		Name: &req.Name,
		Type: &req.Type,
	}
	updateDB, err := s.store.UpdateDevice(&update)
	if err != nil {
		return nil, err
	}
	return updateDB, nil
}

func (s *deviceService) InsertStatusData(data *domain.DeviceMessage) error {
	// valid := s.validateStatusData(data)
	// if !valid {
	// 	return fmt.Errorf("device type and device id mismatch")
	// }

	schedule := store.DeviceSchedule{}
	if scheduleData, ok := data.Schedule["schedule"].(map[string]interface{}); ok {
		if action, ok := scheduleData["action"].(string); ok {
			schedule.Action = action
		}
		if scheduledAt, ok := scheduleData["scheduled_at"].(string); ok {
			parsedTime, err := time.Parse(time.RFC3339, scheduledAt)
			if err == nil {
				schedule.ScheduledAt = parsedTime
			}
		}
		if recurring, ok := scheduleData["recurring"].(string); ok {
			schedule.Recurring = recurring
		}
	}

	insert := store.DeviceStatus{
		DeviceID:  data.DeviceID,
		State:     data.State,
		Schedule:  schedule,
		UpdatedAt: data.Time,
	}

	err := s.store.InsertStatusData(&insert)
	if err != nil {
		return fmt.Errorf("failed to insert device status: %w", err)
	}

	return nil
}
func (s *deviceService) HandleControlMessage(data *domain.DeviceMessage) error {
	// BỔ SUNG: THÊM HÀM ĐỂ KIỂM TRA ACTIVE ĐÃ SYNC

	schedule := store.DeviceSchedule{}
	if scheduleData, ok := data.Schedule["schedule"].(map[string]interface{}); ok {
		if action, ok := scheduleData["action"].(string); ok {
			schedule.Action = action
		}
		if scheduledAt, ok := scheduleData["scheduled_at"].(string); ok {
			parsedTime, err := time.Parse(time.RFC3339, scheduledAt)
			if err == nil {
				schedule.ScheduledAt = parsedTime
			}
		}
		if recurring, ok := scheduleData["recurring"].(string); ok {
			schedule.Recurring = recurring
		}
	}

	insert := store.DeviceStatus{
		DeviceID:  data.DeviceID,
		State:     data.State,
		Schedule:  schedule,
		UpdatedAt: data.Time,
	}

	err := s.store.InsertStatusData(&insert)
	if err != nil {
		return fmt.Errorf("failed to insert device status: %w", err)
	}

	return nil
}

func (s *deviceService) GetDeviceStatus(deviceID int32) (*store.DeviceStatus, error) {
	find := store.FindDevice{ID: &deviceID}
	status, err := s.store.GetStatus(&find)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (s *deviceService) ExecuteCommand(cmd domain.DeviceCommand) error {
	id := cmd.GetDeviceID()
	find := store.FindDevice{ID: &id}
	status, err := s.store.GetStatus(&find)
	if err != nil {
		return err
	}
	if err := cmd.Execute(); err != nil {
		return err
	}
	status.State = cmd.GetState()
	status.UpdatedAt = time.Now()
	statusNEW, err := s.store.UpdateStatus(status)
	if err != nil {
		return err
	}
	history := &store.DeviceHistory{
		DeviceID:  cmd.GetDeviceID(),
		Action:    cmd.GetCommandType(),
		State:     cmd.GetState(),
		By:        "user",
		Timestamp: time.Now(),
	}
	historyNEW, err := s.store.RecordHistory(history)
	if err != nil {
		return err
	}
	return fmt.Errorf("status updated: %+v\nhistory recorded: %+v", statusNEW, historyNEW)
}

func (s *deviceService) ListDeviceHistory(id int32, req *domain.GetDeviceHistoryRequest) ([]*store.DeviceHistory, error) {
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

func (s *deviceService) UpdateStatus(status *store.DeviceStatus) error {
	return nil
}

func (s *deviceService) RecordHistory(history *store.DeviceHistory) error {
	return nil
}

// func (s *deviceService) ListDeviceSchedule(idDevice int32) ([]*store.DeviceSchedule, error) {
// 	find := store.FindDeviceSchedule{DeviceID: &idDevice}
// 	schedules, err := s.store.ListDeviceSchedule(&find)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return schedules, nil
// }
