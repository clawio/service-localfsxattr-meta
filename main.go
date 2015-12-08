package main

import (
	"fmt"
	pb "github.com/clawio/service.localstorexattr.meta/proto/metadata"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"runtime"
	"strconv"
)

const (
	serviceID         = "CLAWIO_LOCALSTOREXATTRMETA"
	dataDirEnvar      = serviceID + "_DATADIR"
	tmpDirEnvar       = serviceID + "_TMPDIR"
	portEnvar         = serviceID + "_PORT"
	propEnvar         = serviceID + "_PROP"
	sharedSecretEnvar = "CLAWIO_SHAREDSECRET"
)

type environ struct {
	dataDir      string
	tmpDir       string
	port         int
	prop         string
	sharedSecret string
}

func getEnviron() (*environ, error) {
	e := &environ{}
	e.dataDir = os.Getenv(dataDirEnvar)
	e.tmpDir = os.Getenv(tmpDirEnvar)
	port, err := strconv.Atoi(os.Getenv(portEnvar))
	if err != nil {
		return nil, err
	}
	e.port = port
	e.prop = os.Getenv(propEnvar)
	e.sharedSecret = os.Getenv(sharedSecretEnvar)
	return e, nil
}
func printEnviron(e *environ) {
	log.Infof("%s=%s", dataDirEnvar, e.dataDir)
	log.Infof("%s=%s", tmpDirEnvar, e.tmpDir)
	log.Infof("%s=%d", portEnvar, e.port)
	log.Infof("%s=%s", propEnvar, e.prop)
	log.Infof("%s=%s", sharedSecretEnvar, "******")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Infof("Service %s started", serviceID)

	env, err := getEnviron()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	printEnviron(env)

	p := &newServerParams{}
	p.dataDir = env.dataDir
	p.tmpDir = env.tmpDir
	p.prop = env.prop
	p.sharedSecret = env.sharedSecret

	srv := newServer(p)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.port))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMetaServer(grpcServer, srv)
	grpcServer.Serve(lis)
}
