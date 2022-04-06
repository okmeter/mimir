package okdb

import "time"

// Config configuration of okdb client
type Config struct {
	Address            string        `yaml:"address"`
	Port               int           `yaml:"port"`
	RefreshTimeout     time.Duration `yaml:"refresh_timeout"`
	RequestTimeout     time.Duration `yaml:"request_timeout"`
	BreakerTimeout     time.Duration `yaml:"breaker_timeout"`
	BreakerMaxFailures uint32        `yaml:"breaker_max_failures"`
	ReceiveMsgSize     int           `yaml:"receive_message_size"`
}

// GetPort of okdb front
func (cfg Config) GetPort() int {
	if cfg.Port > 0 {
		return cfg.Port
	}
	return 8000
}

// GetRefreshTimeout with default value
func (cfg Config) GetRefreshTimeout() time.Duration {
	if cfg.RefreshTimeout > 0 {
		return cfg.RefreshTimeout
	}
	return 10 * time.Second
}

// GetRequestTimeout with default value
func (cfg Config) GetRequestTimeout() time.Duration {
	if cfg.RequestTimeout > 0 {
		return cfg.RequestTimeout
	}
	return 10 * time.Second
}

// GetBreakerTimeout with default value
func (cfg Config) GetBreakerTimeout() time.Duration {
	if cfg.BreakerTimeout > 0 {
		return cfg.BreakerTimeout
	}
	return time.Minute
}

// GetBreakerMaxFailures with default value
func (cfg Config) GetBreakerMaxFailures() uint32 {
	if cfg.BreakerMaxFailures > 0 {
		return cfg.BreakerMaxFailures
	}
	return 5
}

// GetReceiveMsgSize with default value
func (cfg Config) GetReceiveMsgSize() int {
	if cfg.ReceiveMsgSize > 0 {
		return cfg.ReceiveMsgSize
	}
	return 100 << 20
}
