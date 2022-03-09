package main

import "github.com/ethereum/go-ethereum/rpc"

type JsonRequest struct {
	ID      interface{}   `json:"id"`
	Version string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type JsonResponse struct {
	ID      interface{} `json:"id"`
	Version string      `json:"jsonrpc"`
	Error   *JsonError  `json:"error,omitempty"`
	Result  interface{} `json:"result"`
}

type JsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func UnwrapErrorMessage(err error) *JsonError {
	if err == nil {
		return nil
	}
	jerr := &JsonError{Message: err.Error()}
	ec, ok := err.(rpc.Error)
	if ok {
		jerr.Code = ec.ErrorCode()
	}
	de, ok := err.(rpc.DataError)
	if ok {
		jerr.Data = de.ErrorData()
	}
	return jerr
}
