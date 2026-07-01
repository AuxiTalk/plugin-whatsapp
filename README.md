# AuxiTalk WhatsApp Plugin

> WhatsApp plugin for AuxiTalk using whatsmeow.

> Portuguese documentation: [README.pt-BR.md](README.pt-BR.md)

## Overview

This plugin connects AuxiTalk to WhatsApp using QR Code login, receives messages in real time, and sends messages through approved `message.send` calls.

## MVP scope

- QR Code login
- Session persistence
- Receive text messages
- Emit `message.received`
- Send text via `message.send`
- Status events

## Build

```sh
go build -o plugin-whatsapp ./cmd/plugin
```

## Run

```sh
./plugin-whatsapp
```

## Configuration

See `.env.example`.

## Security

Never commit `whatsapp.db` or `.env`.
Session DB gives access to the WhatsApp account.
Sending messages must go through `message.send` capability.
