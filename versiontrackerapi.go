package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"

	pbfc "github.com/brotherlogic/filecopier/proto"
	pb "github.com/brotherlogic/versiontracker/proto"
)

var (
	//Backlog - the print queue
	tracking = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "versiontracker_tracking",
		Help: "The size of the tracking queue",
	}, []string{"server", "versiondate"})
)

//NewVersion a new version is available
func (s *Server) NewVersion(ctx context.Context, req *pb.NewVersionRequest) (*pb.NewVersionResponse, error) {
	s.CtxLog(ctx, fmt.Sprintf("Looking for %v", req.GetVersion().GetJob().GetName()))
	for _, version := range s.tracking {
		if version.GetJob().GetName() == req.GetVersion().GetJob().GetName() {
			s.CtxLog(ctx, fmt.Sprintf("Found %v vs %v", version, req.GetVersion()))
			if version.GetVersionDate() < req.GetVersion().GetVersionDate() {
				tracking.With(prometheus.Labels{"server": req.GetVersion().GetJob().GetName(), "versiondate": fmt.Sprintf("%v", time.Unix(req.GetVersion().GetVersionDate(), 0))}).Set(float64(len(s.tracking)))
				return &pb.NewVersionResponse{}, s.doCopy(ctx, req.GetVersion(), version)
			}
		}
	}
	return &pb.NewVersionResponse{}, nil
}

//NewJob alerts us to a new job running
func (s *Server) NewJob(ctx context.Context, req *pb.NewJobRequest) (*pb.NewJobResponse, error) {
	s.tracking[req.GetVersion().GetJob().GetName()] = req.GetVersion()
	tracking.With(prometheus.Labels{"server": req.GetVersion().GetJob().GetName(), "versiondate": fmt.Sprintf("%v", time.Unix(req.GetVersion().GetVersionDate(), 0))}).Set(float64(len(s.tracking)))
	return &pb.NewJobResponse{}, s.validateVersion(ctx, req.GetVersion().GetJob().GetName())
}

//Callback processes the callback from file copier
func (s *Server) Callback(ctx context.Context, req *pbfc.CallbackRequest) (*pbfc.CallbackResponse, error) {
	version, ok := s.keyTrack[req.GetKey()]
	oldversion, ok2 := s.oldVersion[req.GetKey()]
	s.Log(fmt.Sprintf("CALLBACK %v, %v, %v with %v", req, version, ok, oldversion))
	var err error
	if ok && ok2 {
		//Save the version file alongside the binary
		data, _ := proto.Marshal(version)
		err = ioutil.WriteFile(s.base+version.GetJob().GetName()+".nversion", data, 0644)
		s.Log(fmt.Sprintf("Written the version info (%v) -> %v", err, version))

		if err == nil {
			// Run the shutdown
			if oldversion.GetVersion() == "" {
				s.RaiseIssue("Bad shutdown version", fmt.Sprintf("Shutting down %v on %v -> %v", version.GetJob().GetName(), s.Registry.Identifier, oldversion))
			}
			err = s.slave.shutdown(ctx, oldversion)
			if err != nil {
				s.Log(fmt.Sprintf("SHUTDOWN %v -> %v", time.Now(), err))
			}
		}
	}
	return &pbfc.CallbackResponse{}, err
}
