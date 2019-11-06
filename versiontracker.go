package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbbs "github.com/brotherlogic/buildserver/proto"
	pbgbs "github.com/brotherlogic/gobuildslave/proto"
	pbg "github.com/brotherlogic/goserver/proto"
)

type builder interface {
	getLocal(ctx context.Context, job *pbgbs.Job) (*pbbs.Version, error)
}

type prodBuilder struct {
}

func (p *prodBuilder) getLocal(ctx context.Context, job *pbgbs.Job) (*pbbs.Version, error) {
	file := fmt.Sprintf("/home/simon/gobuild/bin/%v.version", job.Name)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	version := &pbbs.Version{}
	proto.Unmarshal(data, version)

	return version, nil
}

type slave interface {
	list(ctx context.Context, identifier string) ([]*pbgbs.Job, error)
}

type prodSlave struct {
	dial func(server string, identifier string) (*grpc.ClientConn, error)
}

func (p *prodSlave) list(ctx context.Context, identifier string) ([]*pbgbs.Job, error) {
	conn, err := p.dial("gobuildslave", identifier)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbgbs.NewBuildSlaveClient(conn)
	list, err := client.ListJobs(ctx, &pbgbs.ListRequest{})

	if err != nil {
		return nil, err
	}

	jobs := []*pbgbs.Job{}
	for _, job := range list.GetJobs() {
		jobs = append(jobs, job.GetJob())
	}
	return jobs, err
}

//Server main server type
type Server struct {
	*goserver.GoServer
	slave   slave
	builder builder
	jobs    []*pbgbs.Job
	vMap    map[string]*pbbs.Version
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
		jobs:     []*pbgbs.Job{},
	}
	s.slave = &prodSlave{dial: s.DialServer}
	s.builder = &prodBuilder{}
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {

}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "versions", Text: fmt.Sprintf("%v", s.vMap)},
		&pbg.State{Key: "jobs", Value: int64(len(s.jobs))},
	}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	var init = flag.Bool("init", false, "Prep server")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server
	server.RegisterServerV2("versiontracker", false)

	if *init {
		return
	}

	server.RegisterRepeatingTaskNonMaster(server.track, "track", time.Minute*5)
	server.RegisterRepeatingTaskNonMaster(server.buildVersionMap, "build_version_amap", time.Minute*5)

	fmt.Printf("%v", server.Serve())
}
