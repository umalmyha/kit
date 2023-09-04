package bootstrap

import (
	"os"
	"time"
)

type ConfigFunc func(c *Config)

type Config struct {
	shutdownTimeout  time.Duration
	interruptSignals []os.Signal
}

func WithShutdownTimeout(t time.Duration) ConfigFunc {
	return func(c *Config) {
		if t > 0 {
			c.shutdownTimeout = t
		}
	}
}

func WithInterruptSignals(signals ...os.Signal) ConfigFunc {
	return func(c *Config) {
		if len(signals) > 0 {
			c.interruptSignals = append(c.interruptSignals, signals...)
		}
	}
}
