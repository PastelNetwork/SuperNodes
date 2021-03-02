package pastelclient

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/a-ok123/go-psl/internal/common"
	"github.com/a-ok123/go-psl/internal/models"
	"github.com/ybbus/jsonrpc/v2"
)

type PslNode struct {
	client   jsonrpc.RPCClient
	logger   common.Logger
	address  string
	user     string
	password string
}

func (n *PslNode) Init(app *common.Application) {
	n.logger = app.Log
	n.user = app.Cfg.Pastel.User
	n.password = app.Cfg.Pastel.Pwd
	n.address = fmt.Sprintf("http://%s:%d", app.Cfg.Pastel.Host, app.Cfg.Pastel.Port)
}

func (n *PslNode) Connect() {

	n.client = jsonrpc.NewClientWithOpts(n.address,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(n.user+":"+n.password)),
			},
		})
	n.logger.InfoLog.Println("-- Connected to PSL cNode -- ")
}

func (n *PslNode) jsonrpccall(method string, params ...interface{}) (*jsonrpc.RPCResponse, error) {
	if n.client == nil {
		n.Connect()
	}
	response, err := n.client.Call(method, params)
	if err != nil {
		n.logger.ErrorLog.Printf("error making RPC call \"%s\" to %s: [%s]", method, n.address, err)
		return nil, errors.New("RPC call Error")
	}
	if response == nil {
		n.logger.ErrorLog.Printf("empty response on RPC call \"%s\" to %s", method, n.address)
		return nil, errors.New("RPC call Error")
	}
	if response.Error != nil {
		// check response.Error.Code, response.Error.Message and optional response.Error.Data
		n.logger.ErrorLog.Printf("RPC call \"%s\" to %s returns error: %s", method, n.address, response.Error.Message)
		return nil, errors.New("RPC call Error")
	}
	if response.Result == nil {
		n.logger.ErrorLog.Printf("RPC call \"%s\" to %s returns empty result", method, n.address)
		return nil, errors.New("RPC call Error")
	}
	return response, nil
}

func (n *PslNode) jsonrpccallfor(object interface{}, method string, params ...interface{}) error {
	if n.client == nil {
		n.Connect()
	}
	err := n.client.CallFor(&object, method, params)
	if err != nil {
		e, ok := err.(*json.UnmarshalTypeError)
		if ok && e.Value == "string" {
			str := ""
			err := n.client.CallFor(&str, method, params)
			if err == nil {
				return errors.New(str)
			}
		}
		n.logger.ErrorLog.Printf("Error calling RPC method \"%s\"", method)
		return errors.New("RPC call Error")
	}
	if object == nil {
		return errors.New("Nothing found")
	}
	return nil
}

func (n *PslNode) Getblockchaininfo() (*models.Blockchaininfo, error) {
	info := &models.Blockchaininfo{}
	err := n.jsonrpccallfor(info, "getblockchaininfo")
	return info, err
}

func (n *PslNode) ListIDTickets(idtype string) (*[]models.IdTicket, error) {
	tickets := &[]models.IdTicket{}
	err := n.jsonrpccallfor(tickets, "tickets", "list", "id", idtype)
	return tickets, err
}

func (n *PslNode) FindIDTicket(search string) (*models.IdTicket, error) {
	tickets := &models.IdTicket{}
	err := n.jsonrpccallfor(tickets, "tickets", "find", "id", search)
	return tickets, err
}

func (n *PslNode) FindIDTickets(search string) (*[]models.IdTicket, error) {
	tickets := &[]models.IdTicket{}
	err := n.jsonrpccallfor(tickets, "tickets", "find", "id", search)
	return tickets, err
}

func (n *PslNode) ListPastelIDs() (*[]models.PastelID, error) {
	pastelids := &[]models.PastelID{}
	err := n.jsonrpccallfor(pastelids, "pastelid", "list")
	return pastelids, err
}

func (n *PslNode) GetMNRegFee() (int, error) {
	r, err := n.jsonrpccall("storagefee", "getnetworkfee")
	if err != nil {
		return -1, err
	}
	return r.Result.(map[string]interface{})["networkfee"].(int), nil

}
