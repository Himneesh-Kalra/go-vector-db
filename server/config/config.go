package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort      int    `mapstructure:"server_port"`
	DataDir         string `mapstructure:"data_dir"`
	SearchAlgorithm string `mapstructure:"search_algorithm"`
	LSHK            int    `mapstructure:"lsh_k"`
	LSHL            int    `mapstructure:"lsh_l"`
	TCPAddr         string `mapstructure:"tcp_addr"`
}

var AppConfig Config

func LoadConfig() error {
	viper.SetConfigFile("/etc/vecdb/server_config.json")
	viper.SetConfigType("json")

	//defaults
	viper.SetDefault("server_port", 6924)
	viper.SetDefault("data_dir", "data")
	viper.SetDefault("search_algorithm", "brute")
	viper.SetDefault("lsh_k", 10)
	viper.SetDefault("lsh_l", 5)

	viper.SetDefault("tcp_addr", "127.0.0.1:6925")

	//allow env overrides
	viper.AutomaticEnv()

	//read json config
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config: %v", err)
	}

	//unmarshall into struct

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("decode config:  %w", err)
	}

	fmt.Printf("Loaded Config: %+v\n", AppConfig)
	return nil

}
