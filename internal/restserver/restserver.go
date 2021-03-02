package restserver

import (
	"context"
	"fmt"
	"github.com/a-ok123/go-psl/internal/common"
	"github.com/a-ok123/go-psl/internal/pastelclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type RESTServer struct {
	PslNode      pastelclient.PslNode
	config       common.Config
	getHandlers  map[string]interface{}
	postHandlers map[string]interface{}
	wsHandlers   map[string]interface{}
}

func (s *RESTServer) AddGetHandlers(handlers map[string]interface{}) {
	if s.getHandlers == nil {
		s.getHandlers = make(map[string]interface{})
	}
	for key, value := range handlers {
		s.getHandlers[key] = value
	}
}

func (s *RESTServer) Start(app *common.Application) func() error {

	s.config = app.Cfg

	s.AddGetHandlers(map[string]interface{}{
		"/getinfo": s.Getinfo,

		"/tickets/id":     s.GetAllIDTickets,
		"/tickets/id/my":  s.GetMyIDTickets,
		"/tickets/id/:id": s.GetIDTicket,

		"/tickets/mnid":     s.GetAllMNIDTickets,
		"/tickets/mnid/my":  s.GetMyMNIDTickets,
		"/tickets/mnid/:id": s.GetMNIDTicket,

		"/pastelids": s.GetPastelIDs,
	})

	e := echo.New()
	e.Use(middleware.Logger())

	APIRoute := e.Group("/api")
	v1route := APIRoute.Group("/v1")
	for n, h := range s.getHandlers {
		v1route.GET(n, h.(func(echo.Context) error))
	}

	restAddress := fmt.Sprintf("%s:%d", app.Cfg.REST.Host, app.Cfg.REST.Port)

	return app.CreateServer("rest_server",
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

func (s *RESTServer) Getinfo(c echo.Context) error {
	res, err := s.PslNode.Getblockchaininfo()
	if err != nil || res == nil {
		return err
	}
	val := fmt.Sprintf("%d", res.Blocks)
	return c.String(http.StatusOK, val)
}

