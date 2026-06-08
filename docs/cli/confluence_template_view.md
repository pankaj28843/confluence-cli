# confluence template view

Show one Cloud content template

## Synopsis

Show one Confluence Cloud content template by id.

## Examples

confluence template view 12345
  confluence template view 12345 --expand body.storage --json

## Usage

```text
confluence template view <id> [flags]
```

## Options

```text
      --expand strings   Expand value; repeatable or comma-separated
```

## Inherited Options

```text
      --debug             Log HTTP requests to stderr (Authorization header redacted)
      --jq string         Filter JSON output through a jq expression (requires jq on PATH)
      --json              Output as JSON (machine-readable)
      --template string   Render JSON output through a Go text/template
      --timing            Show execution time on stderr
```

## See Also

- [confluence template](confluence_template.md)
