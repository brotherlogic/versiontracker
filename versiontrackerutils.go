package main

import (
	"fmt"
	"io/ioutil"
	"time"

	pbbs "github.com/brotherlogic/buildserver/proto"
	"github.com/golang/protobuf/proto"
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

func (s *Server) doCopy(ctx context.Context, version *pbbs.Version) error {
	// Copy the file over - synchronously
	t := time.Now()
	err := s.copier.copy(ctx, version)
	if err != nil {
		return err
	}

	s.Log(fmt.Sprintf("Copied %v in %v", version.GetJob().GetName(), time.Now().Sub(t)))

	//Save the version file alongside the binary
	data, _ := proto.Marshal(version)
	err = ioutil.WriteFile(s.base+version.GetJob().GetName()+".version", data, 0644)

	if err == nil {
		s.Log(fmt.Sprintf("Requesting shutdown %v", version.GetJob().GetName()))
		s.slave.shutdown(ctx, version.GetJob())
	}

	return err
}
