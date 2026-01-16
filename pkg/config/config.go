package config

import "github.com/spf13/viper"

// BaseConfig 通用基础配置
type BaseConfig struct {
	App  AppConfig  `mapstructure:"app"`
	HTTP HTTPConfig `mapstructure:"http"`
	Log  LogConfig  `mapstructure:"log"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type HTTPConfig struct {
	Port int `mapstructure:"port"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

// Load 加载配置到指定的结构体
func Load(cfg interface{}) error {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	// 支持环境变量
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	return v.Unmarshal(cfg)
}

// LoadBase 加载基础配置
func LoadBase() (*BaseConfig, error) {
	var cfg BaseConfig
	if err := Load(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
