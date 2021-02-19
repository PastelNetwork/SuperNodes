package main

import (
	"context"
	"github.com/a-ok123/go-psl/internal/common"
	"sync"
)

func main() {

	// Connect to cNode

	common.Run( "Pastel Wallet Service",
		"config.yml", "stovacore.log",
		[]func(ctx context.Context, config *common.Config, logger *common.Logger, wg *sync.WaitGroup) func() error{
			// Start RPC Server
			StartJsonRpcServer,
			// Start p2p Listener
			StartP2P,
		})
}
