package domain

const (
	DeviceLight = "light"
	DeviceFan   = "fan"

	DeviceActionTurnON   = "TURN_ON"
	DeviceActionTurnOFF  = "TURN_OFF"
	DeviceActionSetLEVEL = "SET_LEVEL"

	ScheduleRepeatDaily  = "DAILY"
	ScheduleRepeatWeekly = "WEEKLY"
	ScheduleRepeatOnce   = "ONCE"

	MaxlevelDeviceLight = "5"
	MaxlevelDeviceFan   = "3"

	DefaultDeviceLocation = "Living Room"
)

type DeviceCommand interface {
	Execute() error
	GetDeviceID() int32
	GetCommandType() string
	GetState() map[string]interface{}
}

func GetDeviceMaxlevel(deviceType string) string {
	switch deviceType {
	case DeviceLight:
		return MaxlevelDeviceLight
	case DeviceFan:
		return MaxlevelDeviceFan
	default:
		return "Maxlevel Unknown"
	}
}

func IsValidDeviceType(deviceType string) bool {
	switch deviceType {
	case DeviceLight, DeviceFan:
		return true
	default:
		return false
	}
}

func IsValidDeviceAction(action string) bool {
	switch action {
	case DeviceActionTurnON, DeviceActionTurnOFF, DeviceActionSetLEVEL:
		return true
	default:
		return false
	}
}

func IsValidScheduleRepeat(repeat string) bool {
	switch repeat {
	case ScheduleRepeatDaily, ScheduleRepeatWeekly, ScheduleRepeatOnce:
		return true
	default:
		return false
	}
}
