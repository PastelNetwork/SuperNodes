package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/a-ok123/go-psl/internal/common"
	"net/http"
	"sync"
)

/*
	parameters:
		name of art
		number of copies
		price per copy

		artist's pastel id
		address
*/
func RegisterTicket(method common.RpcMethod) ([]byte, error) {

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//err := SendMessage(ctx, "mn1", "Hello to you MN1")
	var err error
	err = nil

	s := fmt.Sprintf(`{"status":"ok", "file": "%s"}`, method.Params[0])

	return []byte(s), err
}

func Getinfo(method common.RpcMethod) ([]byte, error) {
	type info struct {
		Method    string `json:"method"`
		PSLNode   bool   `json:"psl_node"`
		RpcServer bool   `json:"rpc_server"`
	}

	i := info{method.Method, true, true}
	return json.Marshal(i)
}

func StartJsonRpcServer(ctx context.Context, config *common.Config, logger *common.Logger, wg *sync.WaitGroup) func() error {

	rpcServer := common.RpcServer{}
	rpcServer.AddHandler("getinfo", Getinfo)
	rpcServer.AddHandler("regticket", RegisterTicket)

	address := fmt.Sprintf("%s:%d", config.Storage.RpcHost, config.Storage.RpcPort)
	server := rpcServer.InitServer(address)

	return common.CreateServer("rpc_server", ctx, config, logger, wg,
		//startServer
		func(ctx context.Context) error {
			return nil
		},
		//runServer
		func(ctx context.Context) error {
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				return fmt.Errorf("error starting Rest server: %w", err)
			}
			return nil
		},
		//stopServer
		func(ctx context.Context) error {
			return server.Shutdown(ctx)
		})
}
