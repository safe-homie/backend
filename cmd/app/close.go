package app

import "github.com/safe-homie/backend/infrastructure"

func Close(infra infrastructure.AppInfra) {
	infra.Store().Close()
	infra.MQTT().Disconnect()
}
