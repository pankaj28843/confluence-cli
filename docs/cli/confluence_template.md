# confluence template

Cloud content and blueprint templates

## Synopsis

Cloud template read operations.

Confluence Server/Data Center template endpoints are not exposed in the
current official REST OpenAPI, so typed template commands are Cloud-only.

## Examples

confluence template list --limit 10
  confluence template blueprint list --space ENG --json
  confluence template view 12345 --expand body.storage

## Usage

```text
confluence template
```

## Commands

- `blueprint` - Cloud blueprint templates
- `list` - List Cloud content templates
- `view` - Show one Cloud content template

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
