package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/egansoft/breezy/config"
	"github.com/egansoft/breezy/routing"
	"github.com/egansoft/breezy/server"
	"github.com/egansoft/breezy/utils"
)

var helpMsg = `Usage: breezy [OPTIONS] FILE
Run a Breezy server using the routes defined in FILE.

FILE consists of a newline seperated list of routes, which have the forms:
  shell command:
      /my/url/path $ command to run
    Args can be specified with brackets:
      /my/url/[arg1]/path $ command to [arg1] run
    The request body is piped into stdin.
  filesystem root:
      /url/path : relative/filesystem/path
    Static files from the filesystem path are served on the url path, as if the
    filesystem were mounted on the url path.

If multiple routes match a given url, the first one is used.

Options:
  -p, --port     specify the port to run on, with 8080 as the default
  -s, --shell    specify a shell to use, with sh as the default
  -d, --debug    enable debug mode
  -h, --help     display this message and exit
`

func main() {
	filename := parseArgs()
	routes, err := utils.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		usageAndExit()
	}

	router, err := routing.Parse(routes)
	if err != nil {
		fmt.Println(err)
		usageAndExit()
	}

	s := server.New(router)
	s.Start()
}

func parseArgs() string {
	port := flag.Int("port", config.Port, "")
	flag.IntVar(port, "p", config.Port, "")
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

	config.Port = *port
	config.DebugMode = *debug
	config.SetShell(*shell)

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Incorrect number of arguments")
		usageAndExit()
	}

	routesFile := flag.Arg(0)
	return routesFile
}

func helpAndExit() {
	fmt.Printf(helpMsg)
	os.Exit(1)
}

func usageAndExit() {
	fmt.Println("See 'breezy --help'")
	os.Exit(1)
}
