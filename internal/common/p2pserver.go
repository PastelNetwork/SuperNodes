package common

import (
	"context"
)

type P2PServer struct {
}

func (s *P2PServer) Start(app *Application) func() error {

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
