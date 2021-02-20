package common

import (
	"encoding/base64"
	"errors"
	"github.com/ybbus/jsonrpc/v2"
)

type PslNode struct {
	client  jsonrpc.RPCClient
	logger  *Logger
	address string
}

func (n *PslNode) Connect(address string, user string, password string, logger *Logger) {
	n.logger = logger
	n.address = address

	n.logger.InfoLog.Println("-- Connecting to PSL cNode -- ")
	n.client = jsonrpc.NewClientWithOpts(address,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password)),
			},
		})
}

func (n *PslNode) jsonrpccall(method string, params ...interface{}) map[string]interface{} {
	response, err := n.client.Call(method, params)
	if err != nil {
		n.logger.ErrorLog.Printf("Error making RPC call \"%s\" to %s: [%s]", method, n.address, err)
		return nil
	}
	if response == nil {
		n.logger.ErrorLog.Printf("Empty response on RPC call \"%s\" to %s", method, n.address)
		return nil
	}
	if response.Error != nil {
		// check response.Error.Code, response.Error.Message and optional response.Error.Data
		n.logger.ErrorLog.Printf("RPC call \"%s\" to %s returns error: ", method, n.address, response.Error.Message)
		return nil
	}
	return response.Result.(map[string]interface{})
}

func (n *PslNode) Getblockchaininfo() map[string]interface{} {
	return n.jsonrpccall("getblockchaininfo")
}

func (n *PslNode) GetMNRegFee() (int, error) {
	r := n.jsonrpccall("storagefee", "getnetworkfee")
	if r != nil {
		return r["networkfee"].(int), nil
	}
	return -1, errors.New("RPC call Error")
}
