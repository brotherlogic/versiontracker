package main

import (
	"golang.org/x/net/context"

	pbbs "github.com/brotherlogic/buildserver/proto"
)

func (s *Server) track(ctx context.Context) error {
	jobs, err := s.slave.list(ctx, s.Registry.Identifier)
	if err == nil {
		s.jobs = jobs
	}

	return err
}

func (s *Server) buildVersionMap(ctx context.Context) error {
	vMap := make(map[string]*pbbs.Version)

	for _, job := range s.jobs {
		lv, _ := s.builder.getLocal(ctx, job)
		vMap[job.GetName()] = lv
	}

	s.vMap = vMap
	return nil
}
