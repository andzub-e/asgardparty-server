package main

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/constants"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/container"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/validator"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/rng"
	"bitbucket.org/electronicjaw/asgardparty-server/utils"
	"context"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"sync"
	"time"
)

// @title Asgard Party Server
// @version 1.0.1
// @description REST API for Asgard Party Slot.
// @host 51.15.117.238:8086
// @host 0.0.0.0:8086
// @BasePath /

func main() {
	utils.Cache.Init()

	now := time.Now()
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	app := container.Build(ctx, wg)

	rngClient := app.Get(constants.RNGName).(rng.Client)
	utils.PatchRand(rngClient)

	logger := app.Get(constants.LoggerName).(*zap.Logger)
	logger.Info("Starting application...")

	server := app.Get(constants.ServerName).(*http.Server)

	binding.Validator = app.Get(constants.ValidatorName).(*validator.Validator)

	go server.Run()

	zap.S().Infof("Up and running (%s)", time.Since(now))
	zap.S().Infof("Got %s signal. Shutting down...", <-utils.WaitTermSignal())

	if err := server.Shutdown(); err != nil {
		zap.S().Errorf("Error stopping server: %s", err)
	}

	wg.Wait()
	zap.S().Info("Service stopped.")
}
