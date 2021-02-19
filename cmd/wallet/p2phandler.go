package main

import (
	"context"
	"github.com/a-ok123/go-psl/internal/common"
	"sync"
)

func StartP2P(ctx context.Context, config *common.Config, logger *common.Logger, wg *sync.WaitGroup) func() error {

	return common.CreateServer("p2p_node", ctx, config, logger, wg,
		//startServer
		func(ctx context.Context) error {
			return nil
		},
		//runServer
		func(ctx context.Context) error {
			return nil
		},
		//stopServer
		func(ctx context.Context) error {
			return nil
		})
}
