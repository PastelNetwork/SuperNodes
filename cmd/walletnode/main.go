package main

import (
	"context"
	"github.com/pastelnetwork/supernodes/internal/common"
	"github.com/pastelnetwork/supernodes/internal/pastelclient"
	"github.com/pastelnetwork/supernodes/internal/restserver"
)

func main() {

	app := common.NewApplication("Pastel Wallet Node", "config.yml", "wnode.log")

	pslNode := pastelclient.New(&app.Cfg, &app.Log)

	restServer := restserver.New(pslNode, &app.Cfg, &app.Log)

	ticketProc := NewTicketProc(pslNode, &app.Cfg, &app.Log)
	restServer.AddGetHandlers(map[string]interface{}{
		"/ws": ticketProc.RegisterArtTicket,
	})

	//p2pServer := fileserver.P2PServer{}

	app.Run([]func(ctx context.Context, a *common.Application) func() error{
			// Start REST Server
			restServer.Start,
			// Start p2p Listener
			//p2pServer.Start,
		})
}
