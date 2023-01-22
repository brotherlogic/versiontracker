package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbb "github.com/brotherlogic/builder/proto"
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
			s.CtxLog(ctx, fmt.Sprintf("%v is equal %v but %v", job.GetName(), lv.GetVersion() != nv.GetVersion(), time.Now().Sub(time.Unix(nv.GetLastBuildTime(), 0))))
		}
		if err == nil && lv.GetVersion() != nv.GetVersion() && time.Now().Sub(time.Unix(nv.GetLastBuildTime(), 0)) < time.Hour*24 && lv.GetVersionDate() < nv.GetVersionDate() {
			s.CtxLog(ctx, fmt.Sprintf("%v -> %v,%v", job.GetName(), lv.GetVersion(), nv.GetVersion()))
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

	//If we can't read remote, we can't copy
	if err != nil {
		return err
	}

	config, err := s.loadConfig(ctx)
	if err != nil {
		return err
	}

	if nv != nil {
		s.CtxLog(ctx, fmt.Sprintf("Reading remote %v -> %v with %v and %v", name, config, nv, time.Since((time.Unix(nv.GetVersionDate(), 0)))))
	}

	if nv != nil && time.Since(time.Unix(nv.GetVersionDate(), 0)) > time.Hour*24*60 && config.BuildBugs[cv.GetJob().GetName()] == 0 {

		issue, err := s.ImmediateIssue(ctx, "Build needed", fmt.Sprintf("According to %v %v was last built %v", s.Registry.Identifier, cv.GetJob().GetName(), time.Unix(lv.GetVersionDate(), 0)), true)
		if err != nil && status.Convert(err).Code() != codes.ResourceExhausted {
			return err
		}
		config.BuildBugs[cv.GetJob().GetName()] = issue.GetNumber()
		err = s.saveConfig(ctx, config)
		if err != nil {
			return err
		}
	} else if nv != nil && time.Since(time.Unix(nv.GetVersionDate(), 0)) > time.Hour*24*7 {
		go func() {
			ctx, cancel := utils.ManualContext(fmt.Sprintf("versiontracker-%v-rebuild", name), time.Minute*10)
			defer cancel()

			conn, err := utils.LFDialServer(ctx, "builder")
			if err != nil {
				return
			}
			defer conn.Close()

			client := pbb.NewBuildClient(conn)
			client.Refresh(ctx, &pbb.RefreshRequest{Job: name})
		}()
	} else if nv != nil && config.BuildBugs[cv.GetJob().GetName()] != 0 {
		val := config.BuildBugs[cv.GetJob().GetName()]
		if time.Since(time.Unix(nv.GetVersionDate(), 0)) <= time.Hour*24*60 {
			err = s.DeleteIssue(ctx, val)
			if err != nil {
				return err
			}
			delete(config.BuildBugs, cv.GetJob().GetName())
			err = s.saveConfig(ctx, config)
			if err != nil {
				return err
			}
		}
	}

	s.CtxLog(ctx, fmt.Sprintf("Working with %v and %v and %v -> %v, %v", cv, nv, lv, err, err2))

	// Do a remote build if we don't have anything
	if cv == nil && nv == nil {
		s.builder.build(ctx, cv.GetJob())
	}

	// Force a copy if local is wrong
	if err2 != nil {
		if nv != nil {
			return s.doCopy(ctx, nv, lv)
		}
		return fmt.Errorf("Unable to reach remote: %v", err)
	}

	if err == nil {
		s.CtxLog(ctx, fmt.Sprintf("Got %v but %v", lv, nv))
		remote.With(prometheus.Labels{"server": name, "versiondate": fmt.Sprintf("%v", time.Unix(nv.GetVersionDate(), 0))}).Inc()
		if nv.GetVersionDate() > lv.GetVersionDate() || lv.GetVersion() != nv.GetVersion() {
			return s.doCopy(ctx, nv, lv)
		}
	}
	return err
}

func (s *Server) doCopy(ctx context.Context, version, oldversion *pbbs.Version) error {
	s.CtxLog(ctx, fmt.Sprintf("COPYING Versions%v -> %v", version, oldversion))

	// Copy the file over - synchronously
	key := time.Now().UnixNano()
	s.keyTrack[key] = version
	s.oldVersion[key] = oldversion
	return s.copier.copy(ctx, version, key)
}
