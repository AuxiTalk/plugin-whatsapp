package plugin

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/auxitalk/plugin-whatsapp/internal/config"
)

type Runtime struct {
	rpc  *RPC
	logs io.Writer
	cfg  config.Config
}

func NewRuntime(input io.Reader, output io.Writer, logs io.Writer, cfg config.Config) *Runtime {
	r := &Runtime{
		rpc:  NewRPC(input, output),
		logs: logs,
		cfg:  cfg,
	}
	r.registerHandlers()
	return r
}

func (r *Runtime) Listen() error {
	fmt.Fprintf(r.logs, "[plugin-whatsapp] ready db=%s\n", r.cfg.DBPath)
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
		"capabilities":     []string{"message.send"},
	}, nil
}

func (r *Runtime) start(_ json.RawMessage) (any, error) {
	fmt.Fprintln(r.logs, "[plugin-whatsapp] started")
	return map[string]any{"started": true}, nil
}

func (r *Runtime) stop(_ json.RawMessage) (any, error) {
	fmt.Fprintln(r.logs, "[plugin-whatsapp] stopped")
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
	if req.Name != "message.send" {
		return nil, fmt.Errorf("capability not found: %s", req.Name)
	}
	return map[string]any{"sent": true}, nil
}
