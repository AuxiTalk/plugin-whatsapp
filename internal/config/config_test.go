package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("WHATSAPP_DB_PATH", "")
	t.Setenv("WHATSAPP_DEVICE_NAME", "")
	t.Setenv("WHATSAPP_AUTO_RECONNECT", "")
	t.Setenv("WHATSAPP_QR_PRINT", "")
	t.Setenv("WHATSAPP_LOG_LEVEL", "")

	cfg := Load()
	if cfg.DBPath != "./whatsapp.db" {
		t.Fatalf("unexpected db path: %s", cfg.DBPath)
	}
	if cfg.DeviceName != "AuxiTalk" {
		t.Fatalf("unexpected device name: %s", cfg.DeviceName)
	}
	if !cfg.AutoReconnect || !cfg.QRPrint {
		t.Fatalf("expected reconnect and qr print enabled by default: %+v", cfg)
	}
	if cfg.LogLevel != "info" {
		t.Fatalf("unexpected log level: %s", cfg.LogLevel)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("WHATSAPP_DB_PATH", "/tmp/session.db")
	t.Setenv("WHATSAPP_DEVICE_NAME", "TestDevice")
	t.Setenv("WHATSAPP_AUTO_RECONNECT", "false")
	t.Setenv("WHATSAPP_QR_PRINT", "0")
	t.Setenv("WHATSAPP_LOG_LEVEL", "debug")

	cfg := Load()
	if cfg.DBPath != "/tmp/session.db" || cfg.DeviceName != "TestDevice" || cfg.LogLevel != "debug" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if cfg.AutoReconnect || cfg.QRPrint {
		t.Fatalf("expected bool env values to be false: %+v", cfg)
	}
}
