package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Config 只读取进程环境变量，不自动加载 .env，也不内置账号密码。
type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBName     string
	DBUsername string
	DBPassword string
}

func Load() (Config, error) {
	config := Config{
		Port:   envOrDefault("PORT", "8080"),
		DBHost: envOrDefault("DB_HOST", "127.0.0.1"), DBPort: envOrDefault("DB_PORT", "3306"),
		DBName: os.Getenv("DB_NAME"), DBUsername: os.Getenv("DB_USERNAME"), DBPassword: os.Getenv("DB_PASSWORD"),
	}
	for name, value := range map[string]string{
		"DB_NAME": config.DBName, "DB_USERNAME": config.DBUsername, "DB_PASSWORD": config.DBPassword,
	} {
		if strings.TrimSpace(value) == "" {
			return Config{}, fmt.Errorf("必须设置环境变量 %s", name)
		}
	}
	for name, value := range map[string]string{"PORT": config.Port, "DB_PORT": config.DBPort} {
		port, err := strconv.Atoi(value)
		if err != nil || port < 1 || port > 65535 {
			return Config{}, fmt.Errorf("%s 必须是 1～65535 的端口号", name)
		}
	}
	return config, nil
}

func (config Config) DSN() string {
	connection := mysql.NewConfig()
	connection.User = config.DBUsername
	connection.Passwd = config.DBPassword
	connection.Net = "tcp"
	connection.Addr = net.JoinHostPort(config.DBHost, config.DBPort)
	connection.DBName = config.DBName
	connection.ParseTime = true
	connection.Loc = time.Local
	connection.Timeout = 5 * time.Second
	connection.ReadTimeout = 5 * time.Second
	connection.WriteTimeout = 5 * time.Second
	connection.Params = map[string]string{"charset": "utf8mb4"}
	return connection.FormatDSN()
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
