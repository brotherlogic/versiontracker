package main

import (
	"testing"

	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"

	pbd "github.com/brotherlogic/discovery/proto"
	pbgbs "github.com/brotherlogic/gobuildslave/proto"
)

func InitTest() *Server {
	s := Init()
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient("./testing")
	s.Registry = &pbd.RegistryEntry{Identifier: "blah"}
	return s
}

type testSlave struct{}

func (p *testSlave) list(ctx context.Context, identifier string) ([]*pbgbs.Job, error) {
	return []*pbgbs.Job{}, nil
}

func TestBasicPull(t *testing.T) {
	s := InitTest()
	s.track(context.Background())
}
