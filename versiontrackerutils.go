package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pbbs "github.com/brotherlogic/buildserver/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
		if lv.GetVersion() != nv.GetVersion() {
			s.Log(fmt.Sprintf("%v is equal %v but %v", job.GetName(), lv.GetVersion() != nv.GetVersion(), time.Now().Sub(time.Unix(nv.GetLastBuildTime(), 0))))
		}
		if err == nil && lv.GetVersion() != nv.GetVersion() && time.Now().Sub(time.Unix(nv.GetLastBuildTime(), 0)) < time.Hour*24 && lv.GetVersionDate() < nv.GetVersionDate() {
			s.Log(fmt.Sprintf("%v -> %v,%v", job.GetName(), lv.GetVersion(), nv.GetVersion()))
			s.needsCopy[job.GetName()] = nv
		}
	}

	return nil
}

var (
	//Backlog - the print queue
	remote = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "versiontracker_remote",
		Help: "The size of the tracking queue",
	}, []string{"server", "versiondate"})
)

func (s *Server) validateVersion(ctx context.Context, name string) error {
	cv := s.tracking[name]
	nv, err := s.builder.getRemote(ctx, cv.GetJob())
	if err == nil {
		remote.With(prometheus.Labels{"server": name, "versiondate": fmt.Sprintf("%v", time.Unix(nv.GetVersionDate(), 0))}).Inc()
		if nv.GetVersionDate() > cv.GetVersionDate() {
			return s.doCopy(ctx, nv)
		}
	}
	return err
}

func (s *Server) doCopy(ctx context.Context, version *pbbs.Version) error {
	// Copy the file over - synchronously
	key := time.Now().UnixNano()
	s.keyTrack[key] = version
	return s.copier.copy(ctx, version, key)
}
