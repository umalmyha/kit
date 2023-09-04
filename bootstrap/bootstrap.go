package bootstrap

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

var ErrShutdownTimeout = errors.New("shutdown timeout")

type Job func(ctx context.Context) error

type Service interface {
	Start() error
	Stop(ctx context.Context) error
}

type Orchestrator struct {
	services []Service
	jobs     []Job
	cfg      Config
}

func New(opts ...ConfigFunc) *Orchestrator {
	cfg := Config{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	if len(cfg.interruptSignals) == 0 {
		cfg.interruptSignals = []os.Signal{os.Interrupt, os.Kill}
	}

	return &Orchestrator{
		cfg: cfg,
	}
}

func (o *Orchestrator) Service(services ...Service) *Orchestrator {
	o.services = append(o.services, services...)
	return o
}

func (o *Orchestrator) Job(jobs ...Job) *Orchestrator {
	o.jobs = append(o.jobs, jobs...)
	return o
}

func (o *Orchestrator) Start() error {
	notifyCtx, stop := signal.NotifyContext(context.Background(), o.cfg.interruptSignals...)
	defer stop()

	grp, grpCtx := errgroup.WithContext(notifyCtx)

	for i := range o.jobs {
		if job := o.jobs[i]; job != nil {
			grp.Go(func() error {
				return job(grpCtx)
			})
		}
	}

	for i := range o.services {
		svc := o.services[i]
		if svc == nil {
			continue
		}

		errCh := make(chan error, 1)

		grp.Go(func() error {
			if err := svc.Start(); err != nil {
				errCh <- err
			}
			return nil
		})

		grp.Go(func() error {
			select {
			case <-grpCtx.Done():
			case err := <-errCh:
				return err
			}

			shutdownCtx, shutdownCancel := o.shutdownContext()
			defer shutdownCancel()

			return svc.Stop(shutdownCtx)
		})
	}

	return grp.Wait()
}

func (o *Orchestrator) shutdownContext() (context.Context, context.CancelFunc) {
	if o.cfg.shutdownTimeout > 0 {
		return context.WithTimeoutCause(context.Background(), o.cfg.shutdownTimeout, ErrShutdownTimeout)
	}
	return context.WithCancel(context.Background())
}
