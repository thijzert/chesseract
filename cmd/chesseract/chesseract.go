package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/thijzert/chesseract/chesseract"
	plumbing "github.com/thijzert/chesseract/internal/web-plumbing"
)

func main() {
	fmt.Printf("Chesseract version: %s\n", chesseract.PackageVersion)

	var configLocation string

	globalSettings := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	globalSettings.StringVar(&configLocation, "f", "~/.config/chesseract/chesseract.json", "Location of configuration file")
	er := globalSettings.Parse(os.Args[1:])
	if er != nil {
		fmt.Fprintf(os.Stderr, "%v\n", er)
		os.Exit(1)
	}

	conf, er := loadConfig(configLocation)
	if er != nil {
		fmt.Fprintf(os.Stderr, "%v\n", er)
		os.Exit(1)
	}

	args := globalSettings.Args()
	command := ""
	if len(args) > 0 {
		command = args[0]
		args = args[1:]
	}

	var err error = fmt.Errorf("invalid command")
	if command == "" {
		err = consoleGame(&conf, args)
	} else if command == "server" {
		err = apiServer(&conf, args)
	}

	er = saveConfig(conf, configLocation)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		if er != nil {
			fmt.Fprintf(os.Stderr, "error saving config: %v\n", er)
		}
		os.Exit(1)
	}
}

func apiServer(conf *Config, args []string) error {
	var listenPort string
	var storageBackend string

	fs := flag.NewFlagSet(os.Args[0]+" server", flag.ContinueOnError)
	fs.StringVar(&listenPort, "listen", "localhost:36819", "IP and port to listen on")
	fs.StringVar(&storageBackend, "storage", "dory:", "DSN for storage backend")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	log.Printf("Starting server...")

	serverConfig := plumbing.ServerConfig{
		Context:    context.Background(),
		StorageDSN: storageBackend,
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
