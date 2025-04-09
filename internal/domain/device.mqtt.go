package domain

type (
	DeviceControlMessage struct {
		ID     int32  `json:"device_id"`
		Status string `json:"status"`
	}
	DeviceScheduleMessage struct {
		ID     int32  `json:"device_id"`
		Action string `json:"action"`
		Time   string `json:"time"`
		Repeat string `json:"repeat"`
	}
)
