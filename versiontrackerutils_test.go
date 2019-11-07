package main

import (
	"testing"

	pbbs "github.com/brotherlogic/buildserver/proto"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"

	pbd "github.com/brotherlogic/discovery/proto"
	pbgbs "github.com/brotherlogic/gobuildslave/proto"
)

func InitTest() *Server {
	s := Init()
	s.slave = &testSlave{}
	s.builder = &testBuild{}
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient("./testing")
	s.Registry = &pbd.RegistryEntry{Identifier: "blah"}
	return s
}

type testSlave struct{}

func (p *testSlave) list(ctx context.Context, identifier string) ([]*pbgbs.Job, error) {
	return []*pbgbs.Job{}, nil
}

type testBuild struct{}

func (t *testBuild) getRemote(ctx context.Context, job *pbgbs.Job) (*pbbs.Version, error) {
	return &pbbs.Version{Version: "one"}, nil
}

func (t *testBuild) getLocal(ctx context.Context, job *pbgbs.Job) (*pbbs.Version, error) {
	return &pbbs.Version{Version: "two"}, nil
}

func TestBasicPull(t *testing.T) {
	s := InitTest()
	s.track(context.Background())
}

func TestReadLocal(t *testing.T) {
	s := InitTest()
	s.jobs = append(s.jobs, &pbgbs.Job{Name: "what"})
	s.buildVersionMap(context.Background())
}
