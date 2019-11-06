package main

import (
	"fmt"

	"golang.org/x/net/context"
)

func (s *Server) track(ctx context.Context) error {
	jobs, err := s.slave.list(ctx, s.Registry.Identifier)
	if err != nil {
		s.Log(fmt.Sprintf("Found %v jobs", len(jobs)))
	}

	return nil
}
