package whatsapp

import (
	"context"
	"fmt"
	"io"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	client *whatsmeow.Client
	logs   io.Writer
	onQR   func(string)
	onMsg  func(Message)
	onConnect func()
	onDisconnect func()
}

type Message struct {
	ChatID   string
	SenderID string
	Text     string
}

func NewClient(dbPath, deviceName string, logs io.Writer) (*Client, error) {
	container, err := sqlstore.New(context.Background(), "sqlite3", "file:"+dbPath+"?_foreign_keys=on", waLog.Noop)
	if err != nil {
		return nil, err
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return nil, err
	}
	if deviceStore == nil {
		deviceStore = container.NewDevice()
	}

	client := whatsmeow.NewClient(deviceStore, waLog.Noop)
	return &Client{
		client: client,
		logs:   logs,
	}, nil
}

func (c *Client) Connect(ctx context.Context, onQR func(string), onMsg func(Message), onConnect, onDisconnect func()) error {
	c.onQR = onQR
	c.onMsg = onMsg
	c.onConnect = onConnect
	c.onDisconnect = onDisconnect

	c.client.AddEventHandler(c.eventHandler)

	if c.client.Store.ID != nil {
		err := c.client.Connect()
		if err == nil && c.onConnect != nil {
			c.onConnect()
		}
		return err
	}

	qrChan, _ := c.client.GetQRChannel(ctx)
	err := c.client.Connect()
	if err != nil {
		return err
	}

	for evt := range qrChan {
		if evt.Event == "code" && c.onQR != nil {
			c.onQR(evt.Code)
		} else if evt.Event == "success" {
			fmt.Fprintln(c.logs, "[whatsapp] login success")
			if c.onConnect != nil {
				c.onConnect()
			}
		}
	}
	return nil
}

func (c *Client) Disconnect() {
	c.client.Disconnect()
	if c.onDisconnect != nil {
		c.onDisconnect()
	}
}

func (c *Client) SendText(chatJID, text string) (string, error) {
	chat, err := types.ParseJID(chatJID)
	if err != nil {
		return "", err
	}
	resp, err := c.client.SendMessage(context.Background(), chat, &waE2E.Message{
		Conversation: &text,
	})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (c *Client) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Info.IsFromMe {
			return
		}
		if v.Message.GetConversation() != "" {
			c.onMsg(Message{
				ChatID:   v.Info.Chat.String(),
				SenderID: v.Info.Sender.String(),
				Text:     v.Message.GetConversation(),
			})
		}
	case *events.Connected:
		if c.onConnect != nil {
			c.onConnect()
		}
	case *events.Disconnected:
		if c.onDisconnect != nil {
			c.onDisconnect()
		}
	}
}
