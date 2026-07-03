package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/auxitalk/plugin-whatsapp/internal/config"
	"github.com/auxitalk/plugin-whatsapp/internal/whatsapp"
)

type WhatsAppClient interface {
	Connect(context.Context, func(string), func(whatsapp.Message), func(), func()) error
	Disconnect()
	SendText(chatJID, text string) (string, error)
}

type Runtime struct {
	rpc    *RPC
	logs   io.Writer
	cfg    config.Config
	client WhatsAppClient
}

func NewRuntime(input io.Reader, output io.Writer, logs io.Writer, cfg config.Config) (*Runtime, error) {
	client, err := whatsapp.NewClient(cfg.DBPath, cfg.DeviceName, logs)
	if err != nil {
		return nil, err
	}
	return NewRuntimeWithClient(input, output, logs, cfg, client), nil
}

func NewRuntimeWithClient(input io.Reader, output io.Writer, logs io.Writer, cfg config.Config, client WhatsAppClient) *Runtime {
	r := &Runtime{
		rpc:    NewRPC(input, output),
		logs:   logs,
		cfg:    cfg,
		client: client,
	}
	r.registerHandlers()
	return r
}

func (r *Runtime) Listen() error {
	fmt.Fprintf(r.logs, "[plugin-whatsapp] ready db=%s\n", r.cfg.DBPath)

	go func() {
		_ = r.client.Connect(context.Background(),
			func(qr string) {
				_ = r.rpc.request("event.emit", map[string]any{
					"type":    "whatsapp.qr",
					"source":  "whatsapp",
					"payload": map[string]any{"code": qr},
				})
			},
			func(msg whatsapp.Message) {
				_ = r.rpc.request("event.emit", map[string]any{
					"type":   "message.received",
					"source": "whatsapp",
					"payload": map[string]any{
						"chatId":   msg.ChatID,
						"senderId": msg.SenderID,
						"text":     msg.Text,
					},
				})
			},
			func() {
				_ = r.rpc.request("event.emit", map[string]any{
					"type":    "whatsapp.connected",
					"source":  "whatsapp",
					"payload": nil,
				})
			},
			func() {
				_ = r.rpc.request("event.emit", map[string]any{
					"type":    "whatsapp.disconnected",
					"source":  "whatsapp",
					"payload": nil,
				})
			},
		)
	}()

	return r.rpc.Listen()
}

func (r *Runtime) registerHandlers() {
	r.rpc.Handle("plugin.handshake", r.handshake)
	r.rpc.Handle("plugin.start", r.start)
	r.rpc.Handle("plugin.stop", r.stop)
	r.rpc.Handle("plugin.health", r.health)
	r.rpc.Handle("capability.call", r.capabilityCall)
}

func (r *Runtime) handshake(_ json.RawMessage) (any, error) {
	return map[string]any{
		"pluginId":        "whatsapp",
		"protocolVersion": "0.1",
		"capabilities": []string{
			"message.send",
			"whatsapp.request_qr",
			"whatsapp.disconnect",
			"whatsapp.status",
		},
	}, nil
}

func (r *Runtime) start(_ json.RawMessage) (any, error) {
	fmt.Fprintln(r.logs, "[plugin-whatsapp] started")
	return map[string]any{"started": true}, nil
}

func (r *Runtime) stop(_ json.RawMessage) (any, error) {
	fmt.Fprintln(r.logs, "[plugin-whatsapp] stopped")
	r.client.Disconnect()
	return map[string]any{"stopped": true}, nil
}

func (r *Runtime) health(_ json.RawMessage) (any, error) {
	return map[string]any{"ok": true, "pluginId": "whatsapp"}, nil
}

func (r *Runtime) capabilityCall(params json.RawMessage) (any, error) {
	var req struct {
		Name  string          `json:"name"`
		Input json.RawMessage `json:"input"`
	}
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, err
	}

	switch req.Name {
	case "message.send":
		return r.handleSend(req.Input)
	case "whatsapp.request_qr":
		return r.handleRequestQR()
	case "whatsapp.disconnect":
		return r.handleDisconnect()
	case "whatsapp.status":
		return r.handleStatus()
	default:
		return nil, fmt.Errorf("capability not found: %s", req.Name)
	}
}

func (r *Runtime) handleSend(input json.RawMessage) (any, error) {
	var msg struct {
		ChatID string `json:"chatId"`
		Text   string `json:"text"`
	}
	if err := json.Unmarshal(input, &msg); err != nil {
		return nil, err
	}
	msgID, err := r.client.SendText(msg.ChatID, msg.Text)
	if err != nil {
		return nil, err
	}
	return map[string]any{"sent": true, "messageId": msgID}, nil
}

func (r *Runtime) handleRequestQR() (any, error) {
	fmt.Fprintln(r.logs, "[whatsapp] request_qr chamado")
	return map[string]any{"requested": true}, nil
}

func (r *Runtime) handleDisconnect() (any, error) {
	r.client.Disconnect()
	fmt.Fprintln(r.logs, "[whatsapp] disconnect chamado")
	return map[string]any{"disconnected": true}, nil
}

func (r *Runtime) handleStatus() (any, error) {
	return map[string]any{"connected": false}, nil
}
