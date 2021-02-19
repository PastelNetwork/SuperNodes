package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type RpcServer struct {
	funcs map[string]interface{}
}

func (rpcServer *RpcServer) AddHandler(handlerName string, handler func(method RpcMethod) ([]byte, error)) {
	if rpcServer.funcs == nil {
		rpcServer.funcs = make(map[string]interface{})
	}
	rpcServer.funcs[handlerName] = handler
}

func (rpcServer *RpcServer) InitServer(address string) *http.Server {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var rpcMethod RpcMethod
		err := json.NewDecoder(r.Body).Decode(&rpcMethod)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if rpcServer.funcs == nil {
			err := fmt.Errorf("rpcMethod methods are not implemented  - %s", rpcMethod.Method)
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		m := reflect.ValueOf(rpcServer.funcs[strings.ToLower(rpcMethod.Method)])
		//m := reflect.ValueOf(&rpcMethod).MethodByName(strings.Title(strings.ToLower(rpcMethod.Method)))
		if m.Kind() == reflect.Invalid || m.IsNil() || m.IsZero() {
			err := fmt.Errorf("rpcMethod method not implemented  - %s", rpcMethod.Method)
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		in := make([]reflect.Value, 1)
		in[0] = reflect.ValueOf(rpcMethod)

		retValue := m.Call(in)
		if retValue == nil || len(retValue) != 2 ||
			retValue[0].Kind() != reflect.Slice ||
			retValue[1].Kind() != reflect.Interface {
			err := fmt.Errorf("invalid rpcMethod method implementation - %s", rpcMethod.Method)
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		if !retValue[1].IsNil() {
			err := retValue[1].Interface().(error)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		retBytes := retValue[0].Interface().([]byte)
		var buf bytes.Buffer
		err = json.Indent(&buf, retBytes, "", "  ")

		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	})
	return &http.Server{Addr: address, Handler: mux}
}

type RpcMethod struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      string   `json:"id"`
}
