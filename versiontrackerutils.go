package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pbbs "github.com/brotherlogic/buildserver/proto"
	"github.com/brotherlogic/goserver/utils"
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
	lv, err2 := s.builder.getLocal(ctx, cv.GetJob())

	s.Log(fmt.Sprintf("Working with %v and %v and %v -> %v, %v", cv, nv, lv, err, err2))
	// Force a copy if local is wrong
	if err2 != nil {
		if nv != nil {
			return s.doCopy(ctx, nv, lv)
		}
		return fmt.Errorf("Unable to reach remote: %v", err)
	}

	if err == nil {
		s.Log(fmt.Sprintf("Got %v but %v", lv, nv))
		remote.With(prometheus.Labels{"server": name, "versiondate": fmt.Sprintf("%v", time.Unix(nv.GetVersionDate(), 0))}).Inc()
		if nv.GetVersionDate() > lv.GetVersionDate() {
			return s.doCopy(ctx, nv, lv)
		}
	}
	return err
}

func (s *Server) doCopy(ctx context.Context, version, oldversion *pbbs.Version) error {
	s.CtxLog(ctx, fmt.Sprintf("COPYING Versions%v -> %v", version, oldversion))

	if oldversion.GetVersion() == "" {
		key, err := utils.GetContextKey(ctx)
		s.RaiseIssue("Bad format here", fmt.Sprintf("%v, %v -> %v,%v", version, oldversion, key, err))
	}

	// Copy the file over - synchronously
	key := time.Now().UnixNano()
	s.keyTrack[key] = version
	s.oldVersion[key] = oldversion
	return s.copier.copy(ctx, version, key)
}
