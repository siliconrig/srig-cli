# srig-cli

CLI tool for the [siliconrig](https://siliconrig.dev) Hardware-as-a-Service platform.

Flash firmware, open serial terminals, and run CI/CD tests on real MCU boards — all from your terminal.

## Install

Download a binary from [Releases](https://github.com/siliconrig/srig-cli/releases), or build from source:

```bash
go install github.com/siliconrig/srig-cli@latest
```

## Authentication

```bash
export SRIG_API_KEY=key_...
```

Or pass `--api-key` to any command.

## Commands

```
srig status                        Show board availability
srig session create --board TYPE   Start a session
srig session list                  List your sessions
srig session end [ID]              End a session (auto-detects active)
srig flash <firmware.bin>          Flash firmware to the board
srig serial                        Open serial terminal (Ctrl+] to exit)
srig power-cycle                   Power cycle the board
srig whoami                        Show your account info
srig version                       Print CLI version
```

## Flags

| Flag | Env | Description |
|------|-----|-------------|
| `--api-key` | `SRIG_API_KEY` | API key |
| `--base-url` | `SRIG_BASE_URL` | API base URL (default: `https://api.srig.io`) |
| `--json` | | JSON output for scripting/CI |

## CI/CD Usage

```yaml
- name: Flash and test
  env:
    SRIG_API_KEY: ${{ secrets.SRIG_API_KEY }}
  run: |
    srig session create --board esp32-s3
    srig flash firmware.bin
    srig serial --timeout 30s --log output.txt
    srig session end
```
