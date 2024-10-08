package config

import "github.com/kelseyhightower/envconfig"

type ctxKey string

const ConfigCtxKey ctxKey = "config"

type Configuration struct {
	ServerConfig   *ServerConfig   `envconfig:"SERVER"`
	CacheConfig    *ServerConfig   `envconfig:"CACHE"`
	HashcashConfig *HashcashConfig `envconfig:"HASHCASH"`
}
type ServerConfig struct {
	Host string `envconfig:"HOST"`
	Port int    `envconfig:"PORT"`
}

type HashcashConfig struct {
	ZerosCount   int `envconfig:"ZEROS_COUNT"`
	Duration     int `envconfig:"DURATION"`
	MaxIteration int `envconfig:"MAX_ITERATION"`
}

func ParseConfig(prefix ...string) (*Configuration, error) {
	var v Configuration
	p := ""
	if len(prefix) > 0 {
		p = prefix[0]
	}

	if err := envconfig.Process(p, &v); err != nil {
		return nil, err
	}

	return &v, nil
}
