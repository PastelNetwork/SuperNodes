package main

import (
	"github.com/a-ok123/go-psl/internal/common"
	"github.com/a-ok123/go-psl/internal/fileserver"
	"github.com/a-ok123/go-psl/internal/pastelclient"
	"github.com/a-ok123/go-psl/internal/restserver"
)

func main() {

	app := common.Application{}
	app.Init("Pastel Wallet Service", "config.yml", "stovacore.log",)

	pslNode := pastelclient.PslNode{}
	pslNode.Init(&app)

	ticketProc := TicketProc{PslNode: pslNode}
	ticketProc.Init(&app)

	restServer := restserver.RESTServer{PslNode: pslNode}
	restServer.AddGetHandlers(map[string]interface{}{
		"/ws": ticketProc.RegisterArtTicket,
	})

	p2pServer := fileserver.P2PServer{}

	app.Run([]func(a *common.Application) func() error{
			// Start REST Server
			restServer.Start,
			// Start p2p Listener
			p2pServer.Start,
		})
}
