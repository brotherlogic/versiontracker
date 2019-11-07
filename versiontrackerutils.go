package main

import (
	"fmt"

	"golang.org/x/net/context"
)

func (s *Server) track(ctx context.Context) error {
	jobs, err := s.slave.list(ctx, s.Registry.Identifier)
	if err == nil {
		s.jobs = jobs
	}

	return err
}

func (s *Server) buildVersionMap(ctx context.Context) error {
	for _, job := range s.jobs {
		lv, _ := s.builder.getLocal(ctx, job)
		nv, err := s.builder.getRemote(ctx, job)
		if err == nil && lv.GetVersion() != nv.GetVersion() {
			s.Log(fmt.Sprintf("%v -> %v,%v", job.GetName(), lv.GetVersion(), nv.GetVersion()))
			s.needsCopy[job.GetName()] = nv
		}
	}

	return nil
}
