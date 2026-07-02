# AGENTS.md

This repository is the official AuxiTalk WhatsApp plugin.

It connects WhatsApp events to AuxiTalk Core through JSON-RPC 2.0 over line-delimited stdio.

## Required context

Read first:

1. `README.md`
2. `docs/ai-development-guide.md`
3. `plugin.json`
4. `internal/plugin/*`
5. `internal/whatsapp/*`

## Required checks

Before finishing code changes:

```sh
gofmt -w <changed-go-files>
go test ./...
```

## Protocol rules

- stdout is reserved for JSON-RPC only.
- stderr is for logs.
- emit WhatsApp events through `event.emit`.
- request risky sends through `action.request` when appropriate.
- never log QR secrets, session data, or message contents unnecessarily.

## Safety rules

- Do not commit WhatsApp session databases.
- Do not commit phone numbers, chats, contacts, QR data, or tokens.
- Message sending must be explicit and auditable.
- Keep end-to-end tests separated from unit tests.

## Product framing

This plugin is one input/output channel. AuxiTalk must remain multi-channel and not WhatsApp-only.
