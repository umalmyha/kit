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

// Serve starts Service(s) and Job(s) lifecycle
func (o *Orchestrator) Serve() error {
	// we create context to track interrupt signals
	notifyCtx, stop := signal.NotifyContext(context.Background(), o.cfg.interruptSignals...)
	defer stop()

	// shutdown context must be passed immediately to goroutines, so services and jobs can be
	// stopped after some timeout
	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
	defer shutdownCancel()

	//grp, grpCtx := errgroup.WithContext(notifyCtx)

	var (
		wg   sync.WaitGroup
		once sync.Once
	)

	errOnce := func(e error) {
		once.Do(func() {
			stop()
		})
	}

	for i := range o.jobs {
		job := o.jobs[i]
		if job == nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := job(shutdownCtx); err != nil {
				errOnce(err)
			}
		}()
	}

	for i := range o.services {
		svc := o.services[i]
		if svc == nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			go func() {

			}()

			if err := svc.Start(); err != nil {

			}
		}()
	}
}
