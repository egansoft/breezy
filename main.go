package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/egansoft/silly/parser"
	"github.com/egansoft/silly/server"
)

var helpMsg = `Silly Server
silly PORT ROUTES
`

func main() {
	port, routesFile := parseArgs()
	router, err := parser.ParseFile(routesFile)
	if err != nil {
		panic(err)
	}

	s := server.New(port, router)
	s.Start()
}

func parseArgs() (uint, string) {
	if len(os.Args) != 3 {
		helpAndExit()
	}

	portString := os.Args[1]
	portInt, err := strconv.Atoi(portString)
	if err != nil {
		helpAndExit()
	}
	port := uint(portInt)

	routesFile := os.Args[2]
	return port, routesFile
}

func helpAndExit() {
	fmt.Printf(helpMsg)
	os.Exit(1)
}
