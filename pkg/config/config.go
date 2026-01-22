package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// BaseConfig 通用基础配置
type BaseConfig struct {
	App  AppConfig  `mapstructure:"app"`
	HTTP HTTPConfig `mapstructure:"http"`
	Log  LogConfig  `mapstructure:"log"`
	DB   DBConfig   `mapstructure:"db"`
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

type DBConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	SSLMode      string `mapstructure:"sslmode"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxLifetime  int    `mapstructure:"max_lifetime"`
}

// Load 加载配置到指定的结构体
func Load(cfg interface{}, env string) error {
	v := viper.New()

	workDir, err := os.Getwd() // 获取当前工作目录
	if err != nil {
		return err
	}

	// 加载 .env 文件，并将键值对加载到环境变量中
	envPath := filepath.Join(workDir, ".env") // 配置文件路径
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil { // 将 .env 文件中的键值对加载到环境变量中
			fmt.Printf("⚠️  Warning: failed to load .env file: %v\n", err)
		}
	}

	v.SetConfigName("config") // 配置文件名
	v.SetConfigType("yaml")   // 配置文件类型

	if env != "" && env != "dev" {
		v.SetConfigName("config." + env) // 会自动覆盖基础配置
	}

	configPath := filepath.Join(workDir, "configs")

	v.AddConfigPath(configPath) // 配置文件路径
	v.AddConfigPath(workDir)    // 配置文件路径<当前目录>
	v.SetEnvPrefix("APP")       // 环境变量前缀
	v.AutomaticEnv()            // 支持环境变量

	// 设置环境变量替换规则，将 . 替换为 _
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	// 验证：打印加载的配置文件路径
	fmt.Printf("✅ Config file loaded: %s\n", v.ConfigFileUsed())

	// 验证：打印关键配置值
	fmt.Printf("App Env: %s, HTTP Port: %d, Log Level: %s, DB Host: %s, DB Port: %d, DB User: %s, DB Password: %s, DB Name: %s, DB SSL Mode: %s, DB Max Idle Conns: %d, DB Max Open Conns: %d, DB Max Lifetime: %d\n",
		v.GetString("app.env"),
		v.GetInt("http.port"),
		v.GetString("log.level"),
		v.GetString("db.host"),
		v.GetInt("db.port"),
		v.GetString("db.user"),
		v.GetString("db.password"),
		v.GetString("db.dbname"),
		v.GetString("db.sslmode"),
	)

	return v.Unmarshal(cfg)
}
