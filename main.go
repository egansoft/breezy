package main

import (
	"fmt"
	"os"

	"github.com/egansoft/breezy/config"
	"github.com/egansoft/breezy/routing"
	"github.com/egansoft/breezy/server"
)

var helpMsg = `Breezy Server
breezy PORT ROUTES
`

func main() {
	port, routesFile, err := config.ParseArgs()
	if err != nil {
		fmt.Println(err)
		helpAndExit()
	}

	router, err := routing.ParseFile(routesFile)
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
