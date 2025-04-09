package domain

const (
	DeviceLight = "light"
	DeviceFan   = "fan"

	DeviceStatusOn  = "ON"
	DeviceStatusOff = "OFF"

	ScheduleRepeatDaily  = "daily"
	ScheduleRepeatWeekly = "weekly"
	ScheduleDeviceOnce   = "once"

	MaxLevelMode = "5" // not yet

	DefaultDeviceStatus   = "OFF"
	DefaultDeviceLocation = "Living Room"
)

func IsValidDeviceType(deviceType string) bool {
	switch deviceType {
	case DeviceLight, DeviceFan:
		return true
	default:
		return false
	}
}

func IsValidScheduleRepeat(repeat string) bool {
	switch repeat {
	case ScheduleRepeatDaily, ScheduleRepeatWeekly, ScheduleDeviceOnce:
		return true
	default:
		return false
	}
}
