package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/auxitalk/plugin-whatsapp/internal/config"
	"github.com/auxitalk/plugin-whatsapp/internal/whatsapp"
)

type fakeClient struct {
	connectErr   error
	sendErr      error
	sentChatID   string
	sentText     string
	disconnected bool
	onQR         func(string)
	onMsg        func(whatsapp.Message)
	onConnect    func()
	onDisconnect func()
}

func (f *fakeClient) Connect(ctx context.Context, onQR func(string), onMsg func(whatsapp.Message), onConnect, onDisconnect func()) error {
	f.onQR = onQR
	f.onMsg = onMsg
	f.onConnect = onConnect
	f.onDisconnect = onDisconnect
	return f.connectErr
}

func (f *fakeClient) Disconnect() {
	f.disconnected = true
	if f.onDisconnect != nil {
		f.onDisconnect()
	}
}

func (f *fakeClient) SendText(chatJID, text string) (string, error) {
	f.sentChatID = chatJID
	f.sentText = text
	if f.sendErr != nil {
		return "", f.sendErr
	}
	return "msg-1", nil
}

func TestRuntimeHandshakeAndHealth(t *testing.T) {
	runtime := NewRuntimeWithClient(strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{}, config.Config{DBPath: "test.db"}, &fakeClient{})

	result, err := runtime.handshake(nil)
	if err != nil {
		t.Fatalf("handshake: %v", err)
	}
	payload := result.(map[string]any)
	if payload["pluginId"] != "whatsapp" {
		t.Fatalf("unexpected plugin id: %+v", payload)
	}

	health, err := runtime.health(nil)
	if err != nil {
		t.Fatalf("health: %v", err)
	}
	if health.(map[string]any)["ok"] != true {
		t.Fatalf("unexpected health: %+v", health)
	}
}

func TestRuntimeCapabilitySendText(t *testing.T) {
	client := &fakeClient{}
	runtime := NewRuntimeWithClient(strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{}, config.Config{DBPath: "test.db"}, client)
	input, _ := json.Marshal(map[string]any{"chatId": "chat-1", "text": "hello"})
	params, _ := json.Marshal(struct {
		Name  string          `json:"name"`
		Input json.RawMessage `json:"input"`
	}{Name: "message.send", Input: input})

	result, err := runtime.capabilityCall(params)
	if err != nil {
		t.Fatalf("capability call: %v", err)
	}
	if client.sentChatID != "chat-1" || client.sentText != "hello" {
		t.Fatalf("unexpected sent message: chat=%s text=%s", client.sentChatID, client.sentText)
	}
	if result.(map[string]any)["messageId"] != "msg-1" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestRuntimeCapabilitySendTextReturnsError(t *testing.T) {
	runtime := NewRuntimeWithClient(strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{}, config.Config{DBPath: "test.db"}, &fakeClient{sendErr: errors.New("send failed")})
	input, _ := json.Marshal(map[string]any{"chatId": "chat-1", "text": "hello"})
	params, _ := json.Marshal(struct {
		Name  string          `json:"name"`
		Input json.RawMessage `json:"input"`
	}{Name: "message.send", Input: input})

	_, err := runtime.capabilityCall(params)
	if err == nil || !strings.Contains(err.Error(), "send failed") {
		t.Fatalf("expected send error, got %v", err)
	}
}

func TestRuntimeDisconnectCapability(t *testing.T) {
	client := &fakeClient{}
	runtime := NewRuntimeWithClient(strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{}, config.Config{DBPath: "test.db"}, client)
	params, _ := json.Marshal(map[string]any{"name": "whatsapp.disconnect"})

	result, err := runtime.capabilityCall(params)
	if err != nil {
		t.Fatalf("disconnect: %v", err)
	}
	if !client.disconnected || result.(map[string]any)["disconnected"] != true {
		t.Fatalf("expected disconnected result: client=%+v result=%+v", client, result)
	}
}

func TestRuntimeRejectsUnknownCapability(t *testing.T) {
	runtime := NewRuntimeWithClient(strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{}, config.Config{DBPath: "test.db"}, &fakeClient{})
	params, _ := json.Marshal(map[string]any{"name": "whatsapp.unknown"})
	_, err := runtime.capabilityCall(params)
	if err == nil || !strings.Contains(err.Error(), "capability not found") {
		t.Fatalf("expected capability error, got %v", err)
	}
}
