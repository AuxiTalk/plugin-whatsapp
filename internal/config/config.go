package config

import (
	"os"
)

type Config struct {
	DBPath          string
	DeviceName      string
	AutoReconnect   bool
	QRPrint         bool
	LogLevel        string
}

func Load() Config {
	return Config{
		DBPath:        envOrDefault("WHATSAPP_DB_PATH", "./whatsapp.db"),
		DeviceName:    envOrDefault("WHATSAPP_DEVICE_NAME", "AuxiTalk"),
		AutoReconnect: envBool("WHATSAPP_AUTO_RECONNECT", true),
		QRPrint:       envBool("WHATSAPP_QR_PRINT", true),
		LogLevel:      envOrDefault("WHATSAPP_LOG_LEVEL", "info"),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		return v == "true" || v == "1"
	}
	return fallback
}
