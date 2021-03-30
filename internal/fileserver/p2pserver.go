package fileserver

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/supernodes/internal/common"
	"golang.org/x/sync/errgroup"
)

type P2PServer struct {
	config       *common.Config
	logger   	 *common.Logger
	dht			 *DHT
}

func New(cfg *common.Config, log *common.Logger) *P2PServer {
	return &P2PServer{config: cfg, logger: log}
}

func (s *P2PServer) Start(ctx context.Context, app *common.Application) func() error {

	var bootstrapNodes []*NetworkNode

	for _, seed := range app.Cfg.P2P.Seeds {
		if seed.Host != "" || seed.Port != "" {
			bootstrapNode := NewNetworkNode(seed.Host, seed.Port)
			bootstrapNodes = append(bootstrapNodes, bootstrapNode)
		}
	}

	var err error
	s.dht, err = NewDHT(&MemoryStore{}, &Options{
		BootstrapNodes: bootstrapNodes,
		IP:             app.Cfg.P2P.Host,
		Port:           app.Cfg.P2P.Port,
		UseStun:        app.Cfg.P2P.Stun,
	})

	return app.CreateServer(ctx, "p2p_node",
		//initServer
		func(ctx context.Context) error {
			s.logger.InfoLog.Println("p2p_node - Opening socket...")
			if err = s.dht.CreateSocket(); err != nil {
				return fmt.Errorf("p2p_node - error openning Socket for p2p server: %s", err)
			}
			s.logger.InfoLog.Println("p2p_node - Socket opened")
			return nil
		},
		//runServer
		func(ctx context.Context) error {
			eg, ctx := errgroup.WithContext(ctx)

			eg.Go(func() error {
				s.logger.InfoLog.Println("p2p_node is listening on " + s.dht.GetNetworkAddr())
				if err = s.dht.Listen(); err != nil && err.Error() != "closed" {
					return fmt.Errorf("p2p_node - error running p2p server: %s", err)
				}
				return nil
			})
			eg.Go(func() error {
				if len(bootstrapNodes) > 0 {
					s.logger.InfoLog.Println("p2p_node - bootstrapping")
					if err = s.dht.Bootstrap(); err != nil {
						return fmt.Errorf("p2p_node - error bootstrapping p2p server: %s", err)
					}
					s.logger.InfoLog.Println("p2p_node - bootstrapping done")
				}
				return nil
			})

			s.logger.InfoLog.Println("p2p_node started")

			if err := eg.Wait(); err != nil {
				return fmt.Errorf("p2p_node - error in p2p server: %s", err)
			}
			return nil
		},
		//stopServer
		func(ctx context.Context) error {
			if err := s.dht.Disconnect(); err != nil {
				return fmt.Errorf("p2p_node - error stopping p2p server: %s", err)
			}
			return nil
		})
}
