package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/utahta/pythonbrew/log"
	"github.com/utahta/pythonbrew/path"
	"github.com/utahta/pythonbrew/subcmd"
)

func main() {
	if err := run(); err != nil {
		l := log.NewFileLogger()
		l.Errorf("An error has occurred: %v", err)
		l.Debugf("%+v", err)
		l.Printf("See more details: %s", path.Log())
		os.Exit(1)
	}
}

func run() error {
	var (
		showHelp    bool
		showVersion bool
	)
	flag.BoolVar(&showHelp, "h", false, "")
	flag.BoolVar(&showHelp, "help", false, "")
	flag.BoolVar(&showVersion, "v", false, "")
	flag.BoolVar(&showVersion, "version", false, "")

	if len(os.Args) <= 1 {
		return subcmd.NewHelp().Run(nil)
	}

	flag.Parse()
	if showHelp {
		return subcmd.NewHelp().Run(nil)
	}
	if showVersion {
		fmt.Println(subcmd.Version)
		return nil
	}

	c, err := subcmd.Repository().Find(os.Args[1])
	if err != nil {
		return subcmd.NewHelp().Run(nil)
	}
	return c.Run(os.Args[1:])
}
