# confluence settings

Cloud settings reads

## Synopsis

Cloud settings read operations.

Settings and space-settings reads are documented in Confluence Cloud REST API
v1. Server/Data Center does not expose the same settings REST groups in the
current official REST reference, so these typed commands are Cloud-only.

## Examples

confluence settings system-info --json
  confluence settings lookandfeel --space ENG
  confluence settings space ENG

## Usage

```text
confluence settings
```

## Commands

- `lookandfeel` - Show Cloud look-and-feel settings
- `space` - Show Cloud space settings
- `system-info` - Show Cloud system information

## Inherited Options

```text
      --debug             Log HTTP requests to stderr (Authorization header redacted)
      --jq string         Filter JSON output through a jq expression (requires jq on PATH)
      --json              Output as JSON (machine-readable)
      --template string   Render JSON output through a Go text/template
      --timing            Show execution time on stderr
```

## See Also

- [confluence](confluence.md)
