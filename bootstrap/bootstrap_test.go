package bootstrap_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/umalmyha/kit/bootstrap"
)

type service struct {
	stop chan struct{}
}

func newService() *service {
	return &service{stop: make(chan struct{}, 1)}
}

func (s *service) Start() error {
	fmt.Println("starting...")
	<-s.stop
	return nil
}

func (s *service) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done(): // timeout has been raised
		return ctx.Err()
	default:
		time.After(300 * time.Millisecond) // simulate some shutdown stuff...
		fmt.Println("stopping...")
		s.stop <- struct{}{}
	}
	return nil
}

func Example_startupAndShutdown() {
	job := func(ctx context.Context) error {
		fmt.Println("starting...")
		<-ctx.Done()
		return nil
	}

	orc := bootstrap.New(
		bootstrap.WithShutdownTimeout(1*time.Second),
		bootstrap.WithInterruptSignals(os.Interrupt, os.Kill),
	).Job(
		job,
	).Service(
		newService(),
	)

	// send interrupt signal after 500ms
	go func() {
		time.Sleep(500 * time.Millisecond)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT) //nolint:errcheck // just for example
	}()

	if err := orc.Serve(); err != nil {
		fmt.Println("failed to start orchestrator processes: ", err)
	}

	fmt.Println("orchestrator is stopped")

	// Output:
	// starting...
	// starting...
	// stopping...
	// orchestrator is stopped
}

func Example_startupError() {
	job := func(_ context.Context) error {
		return errors.New("job error")
	}

	orc := bootstrap.New(
		bootstrap.WithShutdownTimeout(1*time.Second),
		bootstrap.WithInterruptSignals(os.Interrupt, os.Kill),
	).Job(
		job,
	).Service(
		newService(),
	)

	if err := orc.Serve(); err != nil {
		fmt.Println("failed to serve:", err)
	}

	fmt.Println("orchestrator is stopped")

	// Output:
	// starting...
	// stopping...
	// failed to serve: job error
	// orchestrator is stopped
}
