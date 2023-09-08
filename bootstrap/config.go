package bootstrap

import (
	"os"
	"time"
)

// ConfigFunc is aimed to setup configuration struct
type ConfigFunc func(c *Config)

// Config represents Orchestrator configuration
type Config struct {
	shutdownTimeout  time.Duration
	interruptSignals []os.Signal
}

// WithShutdownTimeout allows to specify timeout in Config. There is no timeout if not specified
func WithShutdownTimeout(t time.Duration) ConfigFunc {
	return func(c *Config) {
		if t > 0 {
			c.shutdownTimeout = t
		}
	}
}

// WithInterruptSignals allows to specify orchestrator interrupt signals. os.Interrupt and os.Kill
// are default signals if not specified
func WithInterruptSignals(signals ...os.Signal) ConfigFunc {
	return func(c *Config) {
		if len(signals) > 0 {
			c.interruptSignals = append(c.interruptSignals, signals...)
		}
	}
}
