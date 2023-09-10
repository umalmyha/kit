package bootstrap

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
)

// ErrShutdownTimeout is raised when server failed to stop services within
// timeout provided in Orchestrator configuration
var ErrShutdownTimeout = errors.New("shutdown timeout")

// Job represents function for handling any endless or running once process.
// Passed context will be cancelled as soon as interrupt signal received
type Job func(ctx context.Context) error

// Service represents any process which is controlled by start and stop.
// On Orchestrator`s startup it tries to call Start function of each Service and
// on receiving interrupt signal Stop function of each service is called forcing
// service to be stopped within startup provided in configuration.
type Service interface {
	Start() error
	Stop(ctx context.Context) error
}

// Orchestrator can collect multiple Job and Service instances and handle theirs lifecycle
type Orchestrator struct {
	services []Service
	jobs     []Job
	cfg      Config
}

// New builds new Orchestrator with provided options
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

// Service accepts services for their lifecycle to be handled
func (o *Orchestrator) Service(services ...Service) *Orchestrator {
	o.services = append(o.services, services...)
	return o
}

// Job accepts jobs for their lifecycle to be handled
func (o *Orchestrator) Job(jobs ...Job) *Orchestrator {
	o.jobs = append(o.jobs, jobs...)
	return o
}

// Serve starts Service(s) and Job(s) lifecycle. The first error is returned if it occurs on startup or shutdown
//
//nolint:cyclop // a lot of variables must be caught within goroutines
func (o *Orchestrator) Serve() (serveErr error) {
	// we create context to track interrupt signals
	notifyCtx, stop := signal.NotifyContext(context.Background(), o.cfg.interruptSignals...)
	defer stop()

	// stop context must be passed immediately to service, so it can be cancelled after shutdown timeout
	stopCtx, stopCancel := context.WithCancel(context.Background())
	defer stopCancel()

	var (
		wg   sync.WaitGroup
		once sync.Once
	)

	// once set error and mark context as done, so setup shutdown process for already started services and jobs
	errOnce := func(err error) {
		once.Do(func() {
			serveErr = err
			stop()
		})
	}

	// each job gets shutdown context immediately
	for i := range o.jobs {
		job := o.jobs[i]
		if job == nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			// job must be running until interrupt signal is received
			if err := job(notifyCtx); err != nil {
				errOnce(err)
			}
		}()
	}

	// each service has to track Start and Stop process in separate goroutine
	for i := range o.services {
		svc := o.services[i]
		if svc == nil {
			continue
		}

		errCh := make(chan error, 1)

		go func() {
			if err := svc.Start(); err != nil {
				// we don't call errOnce here, but pass further to Stop goroutine, so it can decide to
				// call or not Stop function of services
				errCh <- err
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			select {
			case <-notifyCtx.Done():
			case err := <-errCh:
				// error came from Start goroutine, so we start shutdown process for other services and jobs
				// and set an error to return
				errOnce(err)
				return
			}

			if err := svc.Stop(stopCtx); err != nil {
				errOnce(err)
			}
		}()
	}

	// wait for interrupt signal or for the first startup error
	<-notifyCtx.Done()

	// we are waiting for jobs and services to stop within specified timeout if it was
	// specified in config
	done := make(chan struct{}, 1)
	shutdownCtx, shutdownCancel := o.shutdownContext()
	defer shutdownCancel()

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-shutdownCtx.Done():
		errOnce(ErrShutdownTimeout)
	}

	return serveErr
}

func (o *Orchestrator) shutdownContext() (context.Context, context.CancelFunc) {
	if o.cfg.shutdownTimeout > 0 {
		return context.WithTimeoutCause(context.Background(), o.cfg.shutdownTimeout, ErrShutdownTimeout)
	}
	return context.WithCancel(context.Background())
}
