package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	TCPAddr string `mapstructure:"tcp_addr"`
}

var CLIConfig Config

func LoadConfig() error {
	viper.SetConfigFile("/etc/vecdb/cli_config.json")
	viper.SetConfigType("json")

	viper.SetDefault("tcp_addr", "127.0.0.1:6925")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config: %v", err)
	}

	if err := viper.Unmarshal(&CLIConfig); err != nil {
		return fmt.Errorf("decode config: %v", err)
	}

	fmt.Printf("loaded config: %+v\n", &CLIConfig)
	return nil
}
