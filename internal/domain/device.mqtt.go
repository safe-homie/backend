package domain

import "time"

type (
	DeviceMessage struct {
		DeviceID int32                  `json:"device_id"`
		State    map[string]interface{} `json:"state"`
		Schedule map[string]interface{} `json:"schedule"`
		Time     time.Time              `json:"time_stamp"`
		Type     string                 `json:""`
	}

//	DeviceScheduleMessage struct {
//		ID     int32  `json:"device_id"`
//		Action string `json:"action"`
//		Time   string `json:"time"`
//		Repeat string `json:"repeat"`
//	}
)
