package config

import (
	"flag"
	"fmt"
	"strconv"
)

func ParseArgs() (uint, string, error) {
	debug := flag.Bool("d", DebugMode, "enable debug mode")
	shell := flag.String("s", Shell, "set the shell to run")
	flag.Parse()

	DebugMode = *debug
	Shell = *shell

	args := flag.Args()
	if len(args) != 2 {
		err := fmt.Errorf("incorrect number of arguments")
		return 0, "", err
	}

	portString := flag.Arg(0)
	portInt, err := strconv.Atoi(portString)
	if err != nil {
		err = fmt.Errorf("port argument must be an integer")
		return 0, "", err
	}
	port := uint(portInt)

	routesFile := flag.Arg(1)
	return port, routesFile, nil
}
