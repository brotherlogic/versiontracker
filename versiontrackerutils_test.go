package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	pbbs "github.com/brotherlogic/buildserver/proto"
	keystoreclient "github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"

	pbd "github.com/brotherlogic/discovery/proto"
	pbgbs "github.com/brotherlogic/gobuildslave/proto"
)

func InitTest() *Server {
	s := Init()

	s.slave = &testSlave{}
	s.builder = &testBuild{}
	s.copier = &testCopy{}
	s.SkipLog = true
	s.SkipIssue = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient("./testing")
	s.Registry = &pbd.RegistryEntry{Identifier: "blah"}
	s.base = ".tmp/"
	os.Mkdir(".tmp", 0700)

	return s
}

type testCopy struct {
	fail bool
}

func (t *testCopy) copy(ctx context.Context, v *pbbs.Version, k int64) error {
	if t.fail {
		return fmt.Errorf("Built to fail")
	}
	return nil
}

type testSlave struct{}

func (p *testSlave) list(ctx context.Context, identifier string) ([]*pbgbs.Job, error) {
	return []*pbgbs.Job{}, nil
}

func (p *testSlave) listversions(ctx context.Context, identifier string) (string, error) {
	return "", nil
}

func (p *testSlave) shutdown(ctx context.Context, v *pbbs.Version, n string) error {
	return nil
}

type testBuild struct{}

func (t *testBuild) getRemote(ctx context.Context, job *pbgbs.Job) (*pbbs.Version, error) {
	return &pbbs.Version{Version: "one", LastBuildTime: time.Now().Unix(), VersionDate: 10}, nil
}

func (t *testBuild) getLocal(ctx context.Context, job *pbgbs.Job) (*pbbs.Version, error) {
	return &pbbs.Version{Version: "two"}, nil
}

func (t *testBuild) build(ctx context.Context, job *pbgbs.Job) error {
	return nil
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

func TestCopy(t *testing.T) {
	s := InitTest()
	err := s.doCopy(context.Background(), &pbbs.Version{Job: &pbgbs.Job{Name: "Hello"}, Version: "yes", Path: "path", Server: "madeup"}, &pbbs.Version{Job: &pbgbs.Job{Name: "Hello"}, Version: "yes", Path: "path", Server: "madeup"})
	if err != nil {
		t.Errorf("Bad copy: %v", err)
	}
}

func TestCopyBad(t *testing.T) {
	s := InitTest()
	s.copier = &testCopy{fail: true}
	err := s.doCopy(context.Background(),
		&pbbs.Version{
			Job: &pbgbs.Job{
				Name: "Hello",
			},
			Version: "yes",
			Path:    "path",
			Server:  "madeup",
		},
		&pbbs.Version{Job: &pbgbs.Job{Name: "Hello"}, Version: "yes", Path: "path", Server: "madeup"},
	)
	if err == nil {
		t.Errorf("Bad copy did not fail")
	}
}
