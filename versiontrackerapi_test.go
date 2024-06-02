package main

import (
	"testing"
	"time"

	pbbs "github.com/brotherlogic/buildserver/proto"
	pbfc "github.com/brotherlogic/filecopier/proto"
	pb "github.com/brotherlogic/versiontracker/proto"
	"golang.org/x/net/context"
)

func TestNewVersion(t *testing.T) {
	s := InitTest()

	_, err := s.NewVersion(context.Background(), &pb.NewVersionRequest{})
	if err != nil {
		t.Errorf("Bad new version")
	}
}

func TestNewVersionWithUpdate(t *testing.T) {
	s := InitTest()
	s.tracking["what"] = &pbbs.Version{Version: "one", LastBuildTime: time.Now().Unix(), VersionDate: 5}

	_, err := s.NewVersion(context.Background(), &pb.NewVersionRequest{Version: &pbbs.Version{Version: "one", LastBuildTime: time.Now().Unix(), VersionDate: 10}})
	if err != nil {
		t.Errorf("Bad new version")
	}
}

func TestCallback(t *testing.T) {
	s := InitTest()

	s.keyTrack[int64(123)] = &pbbs.Version{}

	_, err := s.Callback(context.Background(), &pbfc.CallbackRequest{Key: int64(123)})

	if err != nil {
		t.Errorf("Bad callback: %v", err)
	}
}
