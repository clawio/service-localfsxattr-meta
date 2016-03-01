package main

import (
	"fmt"
	pb "github.com/clawio/service-localfsxattr-meta/proto/metadata"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"runtime"
	"strconv"
)

const (
	serviceID               = "CLAWIO_LOCALFSXATTR_META"
	dataDirEnvar            = serviceID + "_DATADIR"
	tmpDirEnvar             = serviceID + "_TMPDIR"
	portEnvar               = serviceID + "_PORT"
	logLevelEnvar           = serviceID + "_LOGLEVEL"
	propEnvar               = serviceID + "_PROP"
	propMaxActiveEnvar      = serviceID + "_PROPMAXACTIVE"
	propMaxIdleEnvar        = serviceID + "_PROPMAXIDLE"
	propMaxConcurrencyEnvar = serviceID + "_PROPMAXCONCURRENCY"
	sharedSecretEnvar       = "CLAWIO_SHAREDSECRET"
)

type environ struct {
	dataDir            string
	tmpDir             string
	port               int
	logLevel 	   string
	prop               string
	propMaxActive      int
	propMaxIdle        int
	propMaxConcurrency int
	sharedSecret       string
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

	propMaxActive, err := strconv.Atoi(os.Getenv(propMaxActiveEnvar))
	if err != nil {
		return nil, err
	}
	e.propMaxActive = propMaxActive

	propMaxIdle, err := strconv.Atoi(os.Getenv(propMaxIdleEnvar))
	if err != nil {
		return nil, err
	}
	e.propMaxIdle = propMaxIdle

	propMaxConcurrency, err := strconv.Atoi(os.Getenv(propMaxConcurrencyEnvar))
	if err != nil {
		return nil, err
	}
	e.propMaxConcurrency = propMaxConcurrency

	e.sharedSecret = os.Getenv(sharedSecretEnvar)
	return e, nil
}
func printEnviron(e *environ) {
	log.Infof("%s=%s", dataDirEnvar, e.dataDir)
	log.Infof("%s=%s", tmpDirEnvar, e.tmpDir)
	log.Infof("%s=%d", portEnvar, e.port)
	log.Infof("%s=%s", propEnvar, e.prop)
	log.Infof("%s=%d", propMaxActiveEnvar, e.propMaxActive)
	log.Infof("%s=%d", propMaxIdleEnvar, e.propMaxIdle)
	log.Infof("%s=%d", propMaxConcurrencyEnvar, e.propMaxConcurrency)
	log.Infof("%s=%s", sharedSecretEnvar, "******")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	env, err := getEnviron()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

  	l, err := log.ParseLevel(env.logLevel)
        if err != nil {
                l = log.ErrorLevel
        }
        log.SetLevel(l)

	printEnviron(env)
	log.Infof("Service %s started", serviceID)

	p := &newServerParams{}
	p.dataDir = env.dataDir
	p.tmpDir = env.tmpDir
	p.prop = env.prop
	p.sharedSecret = env.sharedSecret
	p.propMaxActive = env.propMaxActive
	p.propMaxIdle = env.propMaxIdle
	p.propMaxConcurrency = env.propMaxConcurrency

	// Create data and tmp dirs
	if err := os.MkdirAll(p.dataDir, 0644); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if err := os.MkdirAll(p.tmpDir, 0644); err != nil {
		log.Error(err)
		os.Exit(1)
	}

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
