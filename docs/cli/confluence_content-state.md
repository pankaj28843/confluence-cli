# confluence content-state

Cloud content state reads

## Synopsis

Cloud content state read operations.

Content states are documented in Confluence Cloud REST API v1. Server/Data
Center does not expose the same content-state REST group in the current
official REST OpenAPI, so these typed commands are Cloud-only.

## Examples

confluence content-state current 12345
  confluence content-state available 12345 --json
  confluence content-state content ENG --state-id 1 --limit 25

## Usage

```text
confluence content-state
```

## Commands

- `available` - List Cloud states available for content
- `content` - List Cloud content with a given state
- `current` - Show current Cloud content state
- `custom` - List Cloud custom content states
- `settings` - Show Cloud content-state settings for a space
- `space` - List suggested Cloud states for a space

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
