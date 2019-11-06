package main

import (
	"fmt"

	"golang.org/x/net/context"
)

func (s *Server) track(ctx context.Context) error {
	jobs, err := s.slave.list(ctx, s.Registry.Identifier)
	s.Log(fmt.Sprintf("%v and %v", jobs, err))
	if err != nil {
		s.jobs = jobs
	}

	return err
}
