package common

import (
	"context"
	"fmt"
	"github.com/a-ok123/go-psl/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type RESTServer struct {
	getHandlers  map[string]interface{}
	postHandlers map[string]interface{}
	wsHandlers   map[string]interface{}
	psl          PslNode
	config       Config
	logger       Logger
}

func (s *RESTServer) AddGetHandlers(handlers map[string]interface{}) {
	if s.getHandlers == nil {
		s.getHandlers = make(map[string]interface{})
	}
	for key, value := range handlers {
		s.getHandlers[key] = value
	}
}

func (s *RESTServer) Start(app *Application) func() error {

	s.config = app.config

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

	restAddress := fmt.Sprintf("%s:%d", app.config.REST.Host, app.config.REST.Port)

	// Initialise cNode client
	s.psl.Init(app)

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
	res, err := s.psl.Getblockchaininfo()
	if err != nil || res == nil {
		return err
	}
	val := fmt.Sprintf("%d", res.Blocks)
	return c.String(http.StatusOK, val)
}

func (s *RESTServer) getAllIDTickets(c echo.Context, idtype string) error {
	t, err := s.psl.ListIDTickets(idtype)
	if err != nil {
		if err.Error() == "Nothing found" {
			return c.JSON(http.StatusOK, t)
		}
		return err
	}
	return c.JSON(http.StatusOK, t)
}
func (s *RESTServer) getMyIDTickets(c echo.Context, idtype string) error {
	p, err := s.psl.ListPastelIDs()
	if err != nil {
		if err.Error() == "Nothing found" {
			return c.JSON(http.StatusOK, p)
		}
		return err
	}
	t := []models.IdTicket{}
	for _, pid := range *p {
		ticket, err := s.psl.FindIDTicket(pid.PastelID)
		if err != nil {
			if err.Error() == "Key is not found" {
				continue
			}
			return err
		}
		if ticket.Ticket.IDType == idtype {
			t = append(t, *ticket)
		}
	}
	return c.JSON(http.StatusOK, t)
}
func (s *RESTServer) getIDTicket(c echo.Context, idtype string) error {
	id := c.Param("id")
	t, err := s.psl.FindIDTicket(id)
	if err != nil {
		if err.Error() == "Key is not found" {
			return c.JSON(http.StatusOK, "")
		}
		if err.Error() == "Nothing found" {
			return c.JSON(http.StatusOK, "")
		}
		return err
	}
	if t.Ticket.IDType != idtype {
		return c.JSON(http.StatusOK, "")
	}
	return c.JSON(http.StatusOK, t)
}

func (s *RESTServer) GetAllIDTickets(c echo.Context) error {
	return s.getAllIDTickets(c, "personal")
}
func (s *RESTServer) GetMyIDTickets(c echo.Context) error {
	return s.getMyIDTickets(c, "personal")
}
func (s *RESTServer) GetIDTicket(c echo.Context) error {
	return s.getIDTicket(c, "personal")
}
func (s *RESTServer) GetAllMNIDTickets(c echo.Context) error {
	return s.getAllIDTickets(c, "mn")
}
func (s *RESTServer) GetMyMNIDTickets(c echo.Context) error {
	return s.getMyIDTickets(c, "mn")
}
func (s *RESTServer) GetMNIDTicket(c echo.Context) error {
	return s.getIDTicket(c, "mn")
}
func (s *RESTServer) GetPastelIDs(c echo.Context) error {
	p, err := s.psl.ListPastelIDs()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, p)
}

func (s *RESTServer) RegisterTicket(c echo.Context) error {
	return c.String(http.StatusOK, "!")
}
