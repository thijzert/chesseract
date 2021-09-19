package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/thijzert/chesseract/chesseract"
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
		err = consoleLocalMultiplayer(&conf, args)
	} else if command == "server" {
		err = apiServer(&conf, args)
	} else if command == "client" {
		err = consoleGame(&conf, args)
	} else if command == "glclient" {
		err = glGame(&conf, args)
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
