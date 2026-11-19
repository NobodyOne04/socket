# socket

A 2019-2020 university lab: console TCP chat in Go with a session-key protector. Kept as it
was written; changes are docs or small fixes only.

## Commands
- Build: `GO111MODULE=off go build -o messenger .` (relative imports; there is no `go.mod`).
- Server: `./messenger -s -n 2 -p 9090`; client: `./messenger -ip 127.0.0.1 -p 9090`.
- No tests.

## Layout
- `messanger.go`: entry point and flags.
- `server/`, `client/`, `protector/`: one package each (see README).
- `docs/adr/`: decision records (`.adr`), index in `docs/adr/README.md`.

## Conventions
- Read the ADR before changing a recorded decision; commit subjects carry its number.
- Don't modernise (modules, renames) unless asked: it's an archive of the lab.
- The protector is not encryption; never describe it as secure.
