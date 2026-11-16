package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type ReminderConfig struct {
	Interval  time.Duration
	Window    time.Duration
	Workers   int
	QueueSize int
}

func LoadReminderConfig() (ReminderConfig, error) {
	var cfg ReminderConfig
	var err error
	if cfg.Interval, err = envDuration("REMINDER_INTERVAL", 30*time.Second); err != nil {
		return cfg, err
	}
	if cfg.Window, err = envDuration("REMINDER_WINDOW", 10*time.Minute); err != nil {
		return cfg, err
	}
	if cfg.Workers, err = envPositiveInt("REMINDER_WORKERS", 3); err != nil {
		return cfg, err
	}
	if cfg.QueueSize, err = envPositiveInt("REMINDER_QUEUE_SIZE", 100); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func envDuration(key string, def time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration like 30s or 10m, got %q", key, v)
	}
	return d, nil
}

func envPositiveInt(key string, def int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer, got %q", key, v)
	}
	return n, nil
}
