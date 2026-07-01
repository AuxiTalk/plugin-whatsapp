package plugin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Handler func(json.RawMessage) (any, error)

type RPC struct {
	reader   *bufio.Scanner
	writer   io.Writer
	handlers map[string]Handler
}

func NewRPC(input io.Reader, output io.Writer) *RPC {
	return &RPC{
		reader:   bufio.NewScanner(input),
		writer:   output,
		handlers: map[string]Handler{},
	}
}

func (r *RPC) Handle(method string, handler Handler) {
	r.handlers[method] = handler
}

func (r *RPC) Listen() error {
	for r.reader.Scan() {
		line := r.reader.Bytes()
		if len(line) == 0 {
			continue
		}
		if err := r.handleLine(line); err != nil {
			return err
		}
	}
	return r.reader.Err()
}

func (r *RPC) handleLine(line []byte) error {
	var req Request
	if err := json.Unmarshal(line, &req); err != nil {
		return r.write(Response{JSONRPC: "2.0", ID: "", Error: &Error{Code: -32700, Message: "parse error"}})
	}

	if req.JSONRPC != "2.0" || req.Method == "" {
		return r.write(Response{JSONRPC: "2.0", ID: req.ID, Error: &Error{Code: -32600, Message: "invalid request"}})
	}

	handler, ok := r.handlers[req.Method]
	if !ok {
		return r.write(Response{JSONRPC: "2.0", ID: req.ID, Error: &Error{Code: -32601, Message: "method not found"}})
	}

	result, err := handler(req.Params)
	if err != nil {
		return r.write(Response{JSONRPC: "2.0", ID: req.ID, Error: &Error{Code: -32603, Message: err.Error()}})
	}

	if req.ID == "" {
		return nil
	}
	return r.write(Response{JSONRPC: "2.0", ID: req.ID, Result: result})
}

func (r *RPC) write(resp Response) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(r.writer, "%s\n", data)
	return err
}
