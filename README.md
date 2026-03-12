# fbay-cli

CLI tool for the [flashbay](https://flashbay.dev) Hardware-as-a-Service platform.

Flash firmware, open serial terminals, and run CI/CD tests on real MCU boards — all from your terminal.

## Install

Download a binary from [Releases](https://github.com/flashbay-dev/fbay-cli/releases), or build from source:

```bash
go install github.com/flashbay-dev/fbay-cli@latest
```

## Authentication

```bash
export FLASHBAY_API_KEY=key_...
```

Or pass `--api-key` to any command.

## Commands

```
fbay status                        Show board availability
fbay session create --board TYPE   Start a session
fbay session list                  List your sessions
fbay session end [ID]              End a session (auto-detects active)
fbay flash <firmware.bin>          Flash firmware to the board
fbay serial                        Open serial terminal (Ctrl+] to exit)
fbay power-cycle                   Power cycle the board
fbay whoami                        Show your account info
fbay version                       Print CLI version
```

## Flags

| Flag | Env | Description |
|------|-----|-------------|
| `--api-key` | `FLASHBAY_API_KEY` | API key |
| `--base-url` | `FLASHBAY_BASE_URL` | API base URL (default: `https://api.fbay.io`) |
| `--json` | | JSON output for scripting/CI |

## CI/CD Usage

```yaml
- name: Flash and test
  env:
    FLASHBAY_API_KEY: ${{ secrets.FLASHBAY_API_KEY }}
  run: |
    fbay session create --board esp32-s3
    fbay flash firmware.bin
    fbay serial --timeout 30s --log output.txt
    fbay session end
```

## License

Apache 2.0
