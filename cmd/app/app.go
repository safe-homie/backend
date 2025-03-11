package app

import (
	"sync"

	"github.com/safe-homie/backend/internal/transport/mqtt"
	"github.com/safe-homie/backend/internal/transport/rest"
)

func Run() error {
	ctx, infra, srv, err := InitApp()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		restServer := rest.NewServer(ctx, infra, srv)
		if err := restServer.Start(); err != nil {
			ctx.Logger().Error("rest server failed: " + err.Error())
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		mqttServer := mqtt.NewServer(ctx, infra, srv)
		if err := mqttServer.Start(); err != nil {
			ctx.Logger().Error("mqtt server failed: " + err.Error())
		}
	}()
	wg.Wait()
	return nil
}
