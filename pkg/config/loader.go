// 新文件
package config

import (
	"os"
)

// GetEnvironment 获取当前运行环境
func GetEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // 默认开发环境
	}
	return env
}

// LoadBase 加载配置
func LoadBase() (*BaseConfig, error) {
	env := GetEnvironment()
	var cfg BaseConfig

	if err := Load(&cfg, env); err != nil {
		return nil, err
	}

	return &cfg, nil
}
