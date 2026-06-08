# confluence smart-link

Smart Link reads

## Synopsis

Cloud Smart Link read operations.

These commands use documented Cloud v2 embeds endpoints. Create/delete operations
are mutations and remain deferred behind explicit safety gates.

## Examples

confluence smart-link view 12345 --json
  confluence smart-link children 12345 --type page --limit 25

## Usage

```text
confluence smart-link
```

## Commands

- `children` - List direct children of a Cloud smart link
- `view` - Show one Cloud smart link

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
