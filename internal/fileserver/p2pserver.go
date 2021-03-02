package fileserver

import (
	"context"
	"github.com/a-ok123/go-psl/internal/common"
)

type P2PServer struct {
}

func (s *P2PServer) Start(app *common.Application) func() error {

	return app.CreateServer("p2p_node",
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
