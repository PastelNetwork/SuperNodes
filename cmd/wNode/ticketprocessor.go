package main

import (
	"github.com/a-ok123/go-psl/internal/common"
	"github.com/a-ok123/go-psl/internal/pastelclient"
	"github.com/a-ok123/go-psl/internal/restserver"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
)

type TicketProc struct {
	pslNode *pastelclient.PslNode
	config  *common.Config
	logger  *common.Logger
}

type WSMessage struct {
	op ws.OpCode
	msg []byte
}

func NewTicketProc(psl *pastelclient.PslNode, cfg *common.Config, log *common.Logger) *TicketProc {
	return &TicketProc{psl, cfg, log}
}

func (p *TicketProc) RegisterArtTicket(c echo.Context) error {
	conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response().Writer)
	if err != nil {
		return err
	}

	cc := c.(*restserver.RESTServerContext)
	//cc.Jobs.Go( func() error {
	go func() error {
		defer conn.Close()

		p.logger.InfoLog.Println("New Ticket Processor started")

		eg, ctx := errgroup.WithContext(cc.AppCtx)

		messages := make(chan WSMessage)

		eg.Go(func() error {
			p.logger.InfoLog.Println("NTP WS Listener started")
			for {
				p.logger.InfoLog.Println("Waiting for message")
				msg, op, err := wsutil.ReadClientData(conn)
				p.logger.InfoLog.Println("NTP WS Worker exiting - error and signal check")
				if err != nil {
					p.logger.ErrorLog.Printf("Error in New Ticket Processor Listener - %w", err)
					return err
				}
				select {
				case <-ctx.Done():
					p.logger.InfoLog.Println("NTP WS Listener exiting")
					return nil
				case messages <- WSMessage{op, msg}:
					continue
				}
			}
		})
		eg.Go(func() error {
			p.logger.InfoLog.Println("NTP WS Worker started")
			for {
				select {
				case <-ctx.Done():
					p.logger.InfoLog.Println("NTP WS Worker exiting")
					err = wsutil.WriteServerMessage(conn, ws.OpClose, nil)
					return nil
				case msg := <-messages:
					p.logger.InfoLog.Printf("Got message - %s (opcode - %c)", msg.msg, msg.op)
					err = wsutil.WriteServerMessage(conn, msg.op, msg.msg)
					if err != nil {
						p.logger.ErrorLog.Printf("Error in New Ticket Processor Listener - %s", err)
						return err
					}
				}
			}
		})
		if err := eg.Wait(); err != nil {
			p.logger.ErrorLog.Printf("Error in New Ticket Processor - %s", err)
			return nil
		}

		p.logger.InfoLog.Println("New Ticket Processor exiting")
		return nil
	}()
	return nil
}