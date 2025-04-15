package domain

const DefaultDeviceID = "default-device-id"

type Notification struct {
	Title string
	Body  string
	Data  map[string]string
}
