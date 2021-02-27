package main

import (
	"github.com/a-ok123/go-psl/internal/common"
)

func main() {

	app := common.Application{Name: "Pastel Wallet Service"}
	restServer := common.RESTServer{}
	p2pServer := common.P2PServer{}

	app.Run("config.yml", "stovacore.log",
		[]func(a *common.Application) func() error{
			// Start REST Server
			restServer.Start,
			// Start p2p Listener
			p2pServer.Start,
		})
}
