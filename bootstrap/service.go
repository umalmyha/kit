package bootstrap

import "context"

type service struct {
	start func() error
	stop  func(ctx context.Context) error
}

func ToService(start func() error, stop func(ctx context.Context) error) *service {
	return &service{
		start: start,
		stop:  stop,
	}
}

func (s *service) Start() error {
	return s.start()
}

func (s *service) Stop(ctx context.Context) error {
	return s.stop(ctx)
}
