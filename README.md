# fbay-cli

CLI tool for the [Flashbay](https://flashbay.dev) Hardware-as-a-Service platform.

The single source of truth for all Flashbay automation. Used directly by developers and wrapped by CI/CD integrations (GitHub Actions, GitLab CI).

**Status:** Scaffold — commands defined but not yet implemented.

## Install

```bash
go install github.com/flashbay-dev/fbay-cli@latest
```

## Usage

```bash
fbay session create --board esp32-s3
fbay flash firmware.bin
fbay serial
fbay session end
fbay status
```

## CI/CD

- **GitHub Actions:** Use [flashbay-action](https://github.com/flashbay-dev/flashbay-action)
- **GitLab CI:** Include `templates/gitlab-ci.yml` from this repo
