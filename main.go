package main

import (
	"fmt"
	"os"

	"github.com/egansoft/silly/config"
	"github.com/egansoft/silly/parser"
	"github.com/egansoft/silly/server"
)

var helpMsg = `Silly Server
silly PORT ROUTES
`

func main() {
	port, routesFile, err := config.ParseArgs()
	if err != nil {
		fmt.Println(err)
		helpAndExit()
	}

	router, err := parser.ParseFile(routesFile)
	if err != nil {
		panic(err)
	}

	s := server.New(port, router)
	s.Start()
}

func helpAndExit() {
	fmt.Printf(helpMsg)
	os.Exit(1)
}
