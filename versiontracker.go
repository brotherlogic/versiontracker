package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbbs "github.com/brotherlogic/buildserver/proto"
	pbfc "github.com/brotherlogic/filecopier/proto"
	pbgbs "github.com/brotherlogic/gobuildslave/proto"
	pbg "github.com/brotherlogic/goserver/proto"
	pb "github.com/brotherlogic/versiontracker/proto"
)

type copier interface {
	copy(ctx context.Context, v *pbbs.Version, key int64) error
}

type prodCopier struct {
	server func() string
	port   func() int32
	dial   func(ctx context.Context, server, host string) (*grpc.ClientConn, error)
}

func (p *prodCopier) copy(ctx context.Context, v *pbbs.Version, key int64) error {
	conn, err := p.dial(ctx, "filecopier", p.server())
	if err != nil {
		return err
	}
	defer conn.Close()
	copier := pbfc.NewFileCopierServiceClient(conn)
	req := &pbfc.CopyRequest{
		InputFile:    v.GetPath(),
		InputServer:  v.GetServer(),
		OutputFile:   "/home/simon/gobuild/bin/" + v.GetJob().GetName() + ".new",
		OutputServer: p.server(),
		Key:          key,
		Callback:     fmt.Sprintf("%v:%v", p.server(), p.port()),
	}

	cr, err := copier.QueueCopy(ctx, req)
	for err == nil && cr.GetStatus() != pbfc.CopyStatus_COMPLETE {
		time.Sleep(time.Second * 5)
		cr, err = copier.QueueCopy(ctx, req)
	}

	return err
}

type builder interface {
	getLocal(ctx context.Context, job *pbgbs.Job) (*pbbs.Version, error)
	getRemote(ctx context.Context, job *pbgbs.Job) (*pbbs.Version, error)
}

type prodBuilder struct {
	dial   func(ctx context.Context, server string) (*grpc.ClientConn, error)
	server string
}

var (
	//Backlog - the print queue
	remoteReq = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "versiontracker_remotereq",
		Help: "The size of the tracking queue",
	}, []string{"server", "reqstr", "resp"})
)

func (p *prodBuilder) getRemote(ctx context.Context, job *pbgbs.Job) (*pbbs.Version, error) {
	conn, err := p.dial(ctx, "buildserver")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbbs.NewBuildServiceClient(conn)
	req := &pbbs.VersionRequest{JustLatest: true, Job: job, Origin: "versiontracker-" + p.server}
	vers, err := client.GetVersions(ctx, req)
	if err != nil {
		return nil, err
	}

	remoteReq.With(prometheus.Labels{
		"server": job.GetName(),
		"reqstr": fmt.Sprintf("%v", req),
		"resp":   fmt.Sprintf("%v", vers)}).Inc()

	if len(vers.GetVersions()) == 0 {
		return nil, fmt.Errorf("No versions returned")
	}

	return vers.GetVersions()[0], err
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
	shutdown(ctx context.Context, job *pbgbs.Job) error
}

type prodSlave struct {
	dial   func(ctx context.Context, server string, identifier string) (*grpc.ClientConn, error)
	server func() string
}

func (p *prodSlave) list(ctx context.Context, identifier string) ([]*pbgbs.Job, error) {
	conn, err := p.dial(ctx, "gobuildslave", identifier)
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

func (p *prodSlave) shutdown(ctx context.Context, job *pbgbs.Job) error {
	conn, err := p.dial(ctx, job.GetName(), p.server())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pbg.NewGoserverServiceClient(conn)
	_, err = client.Shutdown(ctx, &pbg.ShutdownRequest{})
	return err
}

//Server main server type
type Server struct {
	*goserver.GoServer
	slave     slave
	builder   builder
	copier    copier
	jobs      []*pbgbs.Job
	needsCopy map[string]*pbbs.Version
	base      string
	tracking  map[string]*pbbs.Version
	keyTrack  map[int64]*pbbs.Version
}

func (s *Server) getServerName() string {
	return s.Registry.Identifier
}

func (s *Server) getServerPort() int32 {
	return s.Registry.GetPort()
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer:  &goserver.GoServer{},
		jobs:      []*pbgbs.Job{},
		needsCopy: make(map[string]*pbbs.Version),
		tracking:  make(map[string]*pbbs.Version),
		keyTrack:  make(map[int64]*pbbs.Version),
	}
	s.slave = &prodSlave{dial: s.FDialSpecificServer, server: s.getServerName}
	s.builder = &prodBuilder{dial: s.FDialServer}
	s.copier = &prodCopier{dial: s.FDialSpecificServer, server: s.getServerName, port: s.getServerPort}
	s.base = "/home/simon/gobuild/bin/"
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterVersionTrackerServiceServer(server, s)
	pbfc.RegisterFileCopierCallbackServer(server, s)
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
		&pbg.State{Key: "needs_a_copy", Text: fmt.Sprintf("%v", s.needsCopy)},
	}
}

func (s *Server) runCopy(ctx context.Context) error {
	for key, version := range s.needsCopy {
		err := s.doCopy(ctx, version)
		if err == nil {
			delete(s.needsCopy, key)
		}
		return err
	}

	return nil
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
	err := server.RegisterServerV2("versiontracker", false, true)
	if err != nil {
		return
	}

	if *init {
		return
	}

	server.builder = &prodBuilder{dial: server.FDialServer, server: server.Registry.Identifier}

	fmt.Printf("%v", server.Serve())
}
