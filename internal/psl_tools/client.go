package psl_tools

import (
	"encoding/base64"
	"fmt"
	"github.com/ybbus/jsonrpc/v2"
	"log"
)

type PslNode struct {
	client jsonrpc.RPCClient
}

func (n *PslNode) Connect() {
	log.Println("-- Connecting to PSL cNode -- ")

	n.client = jsonrpc.NewClientWithOpts("http://127.0.0.1:9932",
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("rt"+":"+"rt")),
			},
		})
}

type Networkfee struct {
	Fee int `json:"networkfee"`
}

func (n *PslNode) GetMNRegFee() (int, error) {

	var fee Networkfee
	err := n.client.CallFor(&fee, "storagefee", "getnetworkfee")
	if err != nil {
		return -1, fmt.Errorf("could not send command \"storagefee\" to PSL Node: %s", err)
	}
	return fee.Fee, nil
}
