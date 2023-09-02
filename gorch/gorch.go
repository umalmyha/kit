package gorch

import (
	"context"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"sync"
	"time"
)

type ConfigFunc func(c *Config)

type Config struct {
	ShutdownTimeout  time.Duration
	InterruptSignals []os.Signal
}

func WithShutdownTimeout(t time.Duration) ConfigFunc {
	return func(c *Config) {
		if t > 0 {
			c.ShutdownTimeout = t
		}
	}
}

func WithInterruptSignals(signals ...os.Signal) ConfigFunc {
	return func(c *Config) {
		if len(signals) > 0 {
			c.InterruptSignals = append(c.InterruptSignals, signals...)
		}
	}
}

type Handler func(ctx context.Context) error

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

type Orchestrator struct {
	servers  []Server
	handlers []Handler
	cfg      Config
}

func New(opts ...ConfigFunc) *Orchestrator {
	cfg := Config{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	if len(cfg.InterruptSignals) == 0 {
		cfg.InterruptSignals = []os.Signal{os.Interrupt, os.Kill}
	}

	return &Orchestrator{
		cfg: cfg,
	}
}

func (o *Orchestrator) Server(servers ...Server) *Orchestrator {
	o.servers = append(o.servers, servers...)
	return o
}

func (o *Orchestrator) Handler(handlers ...Handler) *Orchestrator {
	o.handlers = append(o.handlers, handlers...)
	return o
}

func (o *Orchestrator) Start() (startErr error) {
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, o.cfg.InterruptSignals...)
	defer signal.Stop(interruptCh)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}
	errCh := make(chan error, 1)

	once := sync.Once{}
	errOnce := func(err error) {
		once.Do(func() {
			cancel()
			errCh <- err
		})
	}

	for i := range o.handlers {
		handler := o.handlers[i]
		if handler == nil {
			continue
		}

		wg.Add(1)
		go func() {
			if err := handler(ctx); err != nil {
				errOnce(err)
			}
			wg.Done()
		}()
	}

	for i := range o.servers {
		srv := o.servers[i]
		if srv == nil {
			continue
		}

		wg.Add(1)
		go func() {
			go func() {
				<-ctx.Done()

				stopCtx, stopCancel := o.shutdownContext()
				defer stopCancel()

				if err := srv.Stop(stopCtx); err != nil {
					errOnce(err)
				}

				wg.Done()
			}()

			if err := srv.Start(); err != nil {
				errOnce(err)
			}
		}()
	}

	select {
	case <-interruptCh:
		cancel()
	case err := <-errCh:
		startErr = err
	}

	wg.Wait()

	return startErr
}

func (o *Orchestrator) Start1() {
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, o.cfg.InterruptSignals...)
	defer signal.Stop(interruptCh)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grp, grpCtx := errgroup.WithContext(ctx)

	for i := range o.handlers {
		if h := o.handlers[i]; h != nil {
			grp.Go(func() error {
				return h(grpCtx)
			})
		}
	}

	for i := range o.servers {
		if srv := o.servers[i]; srv != nil {
			grp.Go(func() error {
				go func() {
					<-grpCtx.Done()

					stopCtx, stopCancel := o.shutdownContext()
					defer stopCancel()

					if err := srv.Stop(stopCtx); err != nil {
						errOnce(err)
					}

					wg.Done()
				}()

				return srv.Start()
			})
		}
	}
}

func (o *Orchestrator) shutdownContext() (context.Context, context.CancelFunc) {
	if o.cfg.ShutdownTimeout > 0 {
		return context.WithTimeout(context.Background(), o.cfg.ShutdownTimeout)
	}
	return context.WithCancel(context.Background())
}
