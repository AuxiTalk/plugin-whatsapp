# AuxiTalk WhatsApp Plugin

> Plugin de WhatsApp para AuxiTalk usando whatsmeow.

> English documentation: [README.md](README.md)

## Visão geral

Este plugin conecta o AuxiTalk ao WhatsApp usando login por QR Code, recebe mensagens em tempo real e envia mensagens por chamadas aprovadas de `message.send`.

## Escopo MVP

- Login por QR Code
- Persistência de sessão
- Receber mensagens de texto
- Emitir `message.received`
- Enviar texto via `message.send`
- Eventos de status

## Build

```sh
go build -o plugin-whatsapp ./cmd/plugin
```

## Executar

```sh
./plugin-whatsapp
```

## Configuração

Veja `.env.example`.

## Segurança

Nunca commite `whatsapp.db` ou `.env`.
O banco de sessão dá acesso à conta do WhatsApp.
Envio de mensagens deve passar por `message.send`.
