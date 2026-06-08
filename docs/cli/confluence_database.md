# confluence database

Database reads

## Synopsis

Cloud Database read operations.

These commands use documented Cloud v2 databases endpoints. Create/delete operations
are mutations and remain deferred behind explicit safety gates.

## Examples

confluence database view 12345 --json
  confluence database children 12345 --type page --limit 25

## Usage

```text
confluence database
```

## Commands

- `children` - List direct children of a Cloud database
- `view` - Show one Cloud database

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
