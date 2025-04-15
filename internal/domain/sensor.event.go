package domain

type SensorThresholdExceedEvent struct {
	Notify Notification
	Token  string
}
