package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"sync"
)

type RESTServer struct {
	psl         PslNode
	GetHandlers map[string]interface{}
}

func (s *RESTServer) Getinfo(c echo.Context) error {
	res := s.psl.Getblockchaininfo()
	var val string
	if res != nil && res["blocks"] != nil {
		blocks, err := res["blocks"].(json.Number).Int64()
		if err == nil {
			val = fmt.Sprintf("%d", blocks)
		}
	}
	return c.String(http.StatusOK, val)
}

func (s *RESTServer) RegisterTicket(c echo.Context) error {
	return c.String(http.StatusOK, "!")
}

func (s *RESTServer) Start(ctx context.Context, config *Config, logger *Logger, wg *sync.WaitGroup) func() error {

	s.GetHandlers = make(map[string]interface{})
	s.GetHandlers["getinfo"] = s.Getinfo

	e := echo.New()
	for n, h := range s.GetHandlers {
		e.GET(n, h.(func(echo.Context) error))
	}

	restAddress := fmt.Sprintf("%s:%d", config.REST.Host, config.REST.Port)

	// Connect to cNode
	pslAddrress := fmt.Sprintf("http://%s:%d", config.Pastel.Host, config.Pastel.Port)
	s.psl.Connect(pslAddrress, config.Pastel.User, config.Pastel.Pwd, logger)

	return CreateServer("rest_server", ctx, config, logger, wg,
		//startServer
		func(ctx context.Context) error {
			return nil
		},
		//runServer
		func(ctx context.Context) error {
			//if err := http.ListenAndServe(restAddress, mux); err != http.ErrServerClosed {
			if err := e.Start(restAddress); err != http.ErrServerClosed {
				return fmt.Errorf("error starting Rest server: %w", err)
			}
			return nil
		},
		//stopServer
		func(ctx context.Context) error {
			//return server.Shutdown(ctx)
			return nil
		})
}
