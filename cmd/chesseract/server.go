package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	plumbing "github.com/thijzert/chesseract/internal/web-plumbing"

	_ "github.com/thijzert/chesseract/internal/storage/sql"
)

func apiServer(conf *Config, args []string) error {
	logVerbose := false
	var listenPort string
	var storageBackend string

	fs := flag.NewFlagSet(os.Args[0]+" server", flag.ContinueOnError)
	fs.StringVar(&listenPort, "listen", "localhost:36819", "IP and port to listen on")
	fs.StringVar(&storageBackend, "storage", "dory:", "DSN for storage backend")
	fs.BoolVar(&logVerbose, "v", false, "Verbosely log all errors sent to clients")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	log.Printf("Starting server...")

	serverConfig := plumbing.ServerConfig{
		Context:    context.Background(),
		StorageDSN: storageBackend,
	}

	if logVerbose {
		serverConfig.ClientErrorLog = os.Stderr
	}

	s, err := plumbing.New(serverConfig)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s", listenPort)

	var srv http.Server
	srv.Handler = s
	return srv.Serve(ln)
}
