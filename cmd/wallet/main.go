package main

import (
	"context"
	"github.com/a-ok123/go-psl/internal/common"
	"sync"
)

func main() {

	restServer := common.RESTServer{}
	p2pServer := P2PServer{}

	common.Run("Pastel Wallet Service",
		"config.yml", "stovacore.log",
		[]func(ctx context.Context, config *common.Config, logger *common.Logger, wg *sync.WaitGroup) func() error{
			// Start REST Server
			restServer.Start,
			// Start p2p Listener
			p2pServer.Start,
		})
}
