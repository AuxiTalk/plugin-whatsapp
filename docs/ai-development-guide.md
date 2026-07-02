# AI Development Guide

This guide is for AI coding agents working on the AuxiTalk WhatsApp plugin.

## Responsibilities

- Connect WhatsApp events to AuxiTalk Core.
- Emit normalized events such as `whatsapp.qr`, `whatsapp.connected`, `message.received`.
- Support controlled message sending.
- Preserve session data safely.

## Safe workflow

1. Inspect `plugin.json`, `internal/plugin`, and `internal/whatsapp`.
2. Add tests for JSON-RPC handlers and config.
3. Avoid real WhatsApp network tests in default unit tests.
4. Run `gofmt` and `go test ./...`.
5. Commit only when requested.
6. Push only when explicitly requested.

## Sensitive areas

- QR codes.
- SQLite session files.
- Phone numbers and message content.
- Automatic message sending.

Default tests must not require a real WhatsApp account.
