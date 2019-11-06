package main

import (
	"golang.org/x/net/context"
)

func (s *Server) track(ctx context.Context) error {
	jobs, err := s.slave.list(ctx, s.Registry.Identifier)
	if err != nil {
		s.jobs = jobs
	}

	return nil
}
