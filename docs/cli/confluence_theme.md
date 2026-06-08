# confluence theme

Cloud theme reads

## Synopsis

Cloud theme read operations.

Theme reads are documented in Confluence Cloud REST API v1. Server/Data Center
does not expose the same theme REST group in the current official REST
reference, so these typed commands are Cloud-only.

## Examples

confluence theme list --limit 10
  confluence theme global --json
  confluence theme space ENG

## Usage

```text
confluence theme
```

## Commands

- `global` - Show selected Cloud global theme
- `list` - List Cloud themes
- `space` - Show selected Cloud space theme
- `view` - Show one Cloud theme

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
