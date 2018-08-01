package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/egansoft/breezy/config"
	"github.com/egansoft/breezy/routing"
	"github.com/egansoft/breezy/server"
)

var helpMsg = `Usage: breezy [OPTIONS] PORT FILE
Run a Breezy server at the given PORT using the routes defined in FILE.

FILE consists of a newline seperated list of routes, which have the forms:
  shell command:
    /url/[arg1]/with/[arg2]/args $ cmd_name --option arg1 arg2
    which runs the command with the arguments specified in the url
  filesystem root:
    /url/path : relative/filesystem/path
    which serves static files from the system path on that url path

If multiple routes match a given url, the first one is used.

Options:
  -s, --shell    specify a shell to use, with sh as the default
  -d, --debug    enable debug mode
  -h, --help     display this message and exit
`

func main() {
	port, routesFile := parseArgs()
	router, err := routing.ParseFile(routesFile)
	if err != nil {
		fmt.Println(err)
		usageAndExit()
	}

	s := server.New(port, router)
	s.Start()
}

func parseArgs() (uint, string) {
	debug := flag.Bool("debug", config.DebugMode, "")
	flag.BoolVar(debug, "d", config.DebugMode, "")
	shell := flag.String("shell", config.Shell, "")
	flag.StringVar(shell, "s", config.Shell, "")
	help := flag.Bool("help", false, "")
	flag.BoolVar(help, "h", false, "")

	flag.Usage = usageAndExit
	flag.Parse()

	if *help {
		helpAndExit()
	}

	config.DebugMode = *debug
	config.Shell = *shell

	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("Incorrect number of arguments")
		usageAndExit()
	}

	portString := flag.Arg(0)
	portInt, err := strconv.Atoi(portString)
	if err != nil {
		fmt.Println("Port argument must be an integer")
		usageAndExit()
	}
	port := uint(portInt)

	routesFile := flag.Arg(1)
	return port, routesFile
}

func helpAndExit() {
	fmt.Printf(helpMsg)
	os.Exit(1)
}

func usageAndExit() {
	fmt.Println("See 'breezy --help'")
	os.Exit(1)
}
