package config

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/constants"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/history"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/overlord"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/rng"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/tracer"
	"fmt"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var config *Config
var once sync.Once

type EngineConfig struct {
	RTP                      string
	BuildVersion             string // any build info
	TurnOnGambleIfAdmissible bool
	Debug                    bool
	IsCheatsAvailable        bool
	MockRNG                  bool
}

type Config struct {
	ServerConfig    *http.Config
	OverlordConfig  *overlord.Config
	HistoryConfig   *history.Config
	RNGConfig       *rng.Config
	ConstantsConfig *constants.Config
	TracerConfig    *tracer.Config
	EngineConfig    *EngineConfig
}

func New() (*Config, error) {
	once.Do(func() {
		config = &Config{}

		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		serverConfig := viper.Sub("server")
		overlordConfig := viper.Sub("overlord")
		historyConfig := viper.Sub("history")
		validationConfig := viper.Sub("game")
		rngConfig := viper.Sub("rng")
		tracerConfig := viper.Sub("tracer")
		engineConfig := viper.Sub("engine")

		if err := parseSubConfig(serverConfig, &config.ServerConfig); err != nil {
			panic(err)
		}

		if err := parseSubConfig(validationConfig, &config.ConstantsConfig); err != nil {
			panic(err)
		}

		if err := parseSubConfig(overlordConfig, &config.OverlordConfig); err != nil {
			panic(err)
		}

		if err := parseSubConfig(historyConfig, &config.HistoryConfig); err != nil {
			panic(err)
		}

		if err := parseSubConfig(rngConfig, &config.RNGConfig); err != nil {
			panic(err)
		}

		if tracerConfig != nil {
			if err := tracerConfig.Unmarshal(&config.TracerConfig); err != nil {
				panic(err)
			}
		} else {
			config.TracerConfig = &tracer.Config{Disabled: true}
		}

		if err := parseSubConfig(engineConfig, &config.EngineConfig); err != nil {
			panic(err)
		}

		config.ServerConfig.MaxProcessingTime *= time.Millisecond
	})

	return config, nil
}

func parseSubConfig[T any](subConfig *viper.Viper, parseTo *T) error {
	if subConfig == nil {
		return fmt.Errorf("can not read %T config: subconfig is nil", parseTo)
	}

	if err := subConfig.Unmarshal(parseTo); err != nil {
		return err
	}

	return nil
}
