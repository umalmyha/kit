package bootstrap

import "context"

type service struct {
	start func() error
	stop  func(ctx context.Context) error
}

// ToService accepts start and stop functions and creates Service based on them
//
//nolint:revive // function is used to pass directly to Orchestrator's Service function
func ToService(start func() error, stop func(ctx context.Context) error) *service {
	return &service{
		start: start,
		stop:  stop,
	}
}

// Start starts process specified by start function on construction
func (s *service) Start() error {
	return s.start()
}

// Stop stops process specified by stop function on construction
func (s *service) Stop(ctx context.Context) error {
	return s.stop(ctx)
}
