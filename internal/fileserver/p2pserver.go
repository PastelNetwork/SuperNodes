package fileserver

import (
	"context"
	"github.com/a-ok123/go-psl/internal/common"
)

type P2PServer struct {
}

func (s *P2PServer) Start(ctx context.Context, app *common.Application) func() error {

	return app.CreateServer(ctx, "p2p_node",
		//initServer
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
