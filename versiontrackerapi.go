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
	for _, version := range s.tracking {
		if version.GetJob().GetName() == req.GetVersion().GetJob().GetName() {
			s.Log(fmt.Sprintf("Found %v vs %v", version, req.GetVersion()))
			if version.GetVersionDate() < req.GetVersion().GetVersionDate() {
				s.tracking[req.GetVersion().GetJob().GetName()] = req.GetVersion()
				tracking.With(prometheus.Labels{"server": req.GetVersion().GetJob().GetName(), "versiondate": fmt.Sprintf("%v", time.Unix(req.GetVersion().GetVersionDate(), 0))}).Set(float64(len(s.tracking)))
				return &pb.NewVersionResponse{}, s.doCopy(ctx, req.GetVersion())
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
	var err error
	if ok {
		//Save the version file alongside the binary
		data, _ := proto.Marshal(version)
		err = ioutil.WriteFile(s.base+version.GetJob().GetName()+".version", data, 0644)

		if err == nil {
			s.Log(fmt.Sprintf("Requesting shutdown %v", version.GetJob().GetName()))
			s.slave.shutdown(ctx, version.GetJob())
		}
	}
	return &pbfc.CallbackResponse{}, err
}
