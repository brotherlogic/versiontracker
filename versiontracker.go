package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
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
	"github.com/brotherlogic/goserver/utils"
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

	_, err = copier.QueueCopy(ctx, req)
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
		"resp":   fmt.Sprintf("%v-%v", len(vers.GetVersions()), vers)}).Inc()

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

	res, err := exec.Command("md5sum", fmt.Sprintf("/home/simon/gobuild/bin/%v", job.GetName())).Output()
	if err != nil {
		return nil, err
	}
	elems := strings.Fields(string(res))
	if version.GetVersion() != elems[0] {
		return nil, fmt.Errorf("Mismatch of versions")
	}

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

var shutdowns = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "versiontracker_shutdowns",
	Help: "Shutdown attempts",
}, []string{"error"})

func (p *prodSlave) shutdown(ctx context.Context, job *pbgbs.Job) error {
	conn, err := p.dial(ctx, job.GetName(), p.server())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pbg.NewGoserverServiceClient(conn)
	_, err = client.Shutdown(ctx, &pbg.ShutdownRequest{})
	shutdowns.With(prometheus.Labels{"error": fmt.Sprintf("%v", err)}).Inc()
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
	server.DiskLog = true
	if err != nil {
		return
	}

	if *init {
		return
	}

	server.builder = &prodBuilder{dial: server.FDialServer, server: server.Registry.Identifier}

	go func() {
		ctx, cancel := utils.ManualContext("versiontrack", time.Minute)
		defer cancel()
		jobs, err := server.slave.list(ctx, server.Registry.GetIdentifier())
		if err != nil {
			log.Fatalf("Cannot reach master: %v", err)
		}

		//Add in the slave itself
		job := &pbgbs.Job{Name: "gobuildslave"}
		lvgbs, err := server.builder.getLocal(ctx, job)
		if err == nil {
			_, err1 := server.NewJob(ctx, &pb.NewJobRequest{Version: lvgbs})
			if err1 != nil {
				server.Log(fmt.Sprintf("Error tracking gbs: %v", err))
			}
		} else {
			_, err1 := server.NewJob(ctx, &pb.NewJobRequest{Version: &pbbs.Version{Job: job}})
			if err1 != nil {
				server.Log(fmt.Sprintf("Error tracking gbs: %v", err))
			}
		}

		server.Log(fmt.Sprintf("Working on %v jobs", len(jobs)))
		time.Sleep(time.Second * 2)
		for _, j := range jobs {
			server.Log(fmt.Sprintf("Working on %v", j))
			ctx, cancel2 := utils.ManualContext("versiontrack", time.Minute)
			defer cancel2()

			lv, err := server.builder.getLocal(ctx, j)
			if err == nil {
				ctx, cancel3 := utils.ManualContext("versiontrack", time.Minute)
				defer cancel3()

				_, err2 := server.NewJob(ctx, &pb.NewJobRequest{Version: lv})
				if err2 != nil {
					server.Log(fmt.Sprintf("Error on new with local job (%v): %v", lv, err2))
				}
			} else {
				ctx, cancel3 := utils.ManualContext("versiontrack", time.Minute)
				defer cancel3()

				_, err2 := server.NewJob(ctx, &pb.NewJobRequest{Version: &pbbs.Version{Job: j}})
				if err2 != nil {
					server.Log(fmt.Sprintf("Error on new job (%v): %v", j, err2))
				}

			}
		}

		server.Log("All jobs processed!")
	}()

	// Prep for shutdown tracking
	err = os.Mkdir("/media/scratch/versiontracker-shutdown", 0777)
	if err != nil {
		server.RaiseIssue(fmt.Sprintf("Tracking dir failure for %v", server.Registry.Identifier), fmt.Sprintf("Dir creation failed: %v", err))
	}

	fmt.Printf("%v", server.Serve())
}
