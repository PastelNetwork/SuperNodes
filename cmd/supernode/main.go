package main

import (
	"context"
	"github.com/pastelnetwork/supernodes/internal/common"
	"github.com/pastelnetwork/supernodes/internal/fileserver"
	"github.com/pastelnetwork/supernodes/internal/pastelclient"
	"github.com/pastelnetwork/supernodes/internal/restserver"
)

func main() {
	app := common.NewApplication("Pastel Super Node", "config.yml", "sNode.log")

	pslNode := pastelclient.New(&app.Cfg, &app.Log)
	restServer := restserver.New(pslNode, &app.Cfg, &app.Log)

	p2pServer := fileserver.New(&app.Cfg, &app.Log)

	app.Run([]func(ctx context.Context, a *common.Application) func() error{
		// Start REST Server
		restServer.Start,
		// Start p2p Listener
		p2pServer.Start,
	})

	//id, err := dht.Store([]byte(input[1]))
	//data, exists, err := dht.Get(input[1])
	//nodes := dht.NumNodes()
	//self := dht.GetSelfID()
	//addr := dht.GetNetworkAddr()
}
