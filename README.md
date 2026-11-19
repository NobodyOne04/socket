# socket

A console TCP chat in Go (university lab 2, 2019-2020): one server, several clients, and a
**session protector** that checks a rotating key on every message.

## Run

It uses GOPATH-era relative imports (`./server`, `./client`), so build it without modules:

```sh
GO111MODULE=off go build -o messenger .
./messenger -s -n 2 -h localhost -p 9090      # server, accepts 2 clients
./messenger -ip 127.0.0.1 -p 9090             # client: enter a login, then type messages
```

| Flag | Side | Default | Meaning |
|---|---|---|---|
| `-s` | server | `false` | start as the server |
| `-n` | server | `2` | how many clients to accept |
| `-h` | server | `localhost` | host to listen on |
| `-p` | both | `9090` | port |
| `-ip` | client | `127.0.0.1` | server address |

## How it works

- `messanger.go`: flags; starts the server or the client in a goroutine.
- `server/server.go`: accepts up to `-n` clients, asks for a login, then relays each message
  to every other client; a goroutine per client reads, one goroutine broadcasts.
- `client/client.go`: sends `hash\nkey` on connect, then `key\nlogin:message` per line.
- `protector/protector.go`: derives the next session key from the current one and the hash.
  Both sides advance the key on every message and disconnect on a mismatch.

The protector is a teaching exercise, not security: the hash travels in plaintext. The same
idea was rebuilt properly in Rust in [TCP-messanger](https://github.com/NobodyOne04/TCP-messanger).

Decisions: [docs/adr](docs/adr/README.md).
