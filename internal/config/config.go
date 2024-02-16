package config

import "time"

func Load() *Config {
	// TODO: add a config loading mechanism like viper
	return &Config{
		Clients: ClientsConfig{
			DeliveryTimeEstimator: DeliveryTimeEstimatorConfig{Timeout: 200 * time.Millisecond},
			DelayCheckQueue:       DelayCheckQueueConfig{Timeout: 200 * time.Millisecond},
		},
	}
}

type Config struct {
	Clients ClientsConfig
}

type ClientsConfig struct {
	DeliveryTimeEstimator DeliveryTimeEstimatorConfig
	DelayCheckQueue       DelayCheckQueueConfig
}

type DeliveryTimeEstimatorConfig struct {
	Timeout time.Duration
}
type DelayCheckQueueConfig struct {
	Timeout time.Duration
}
