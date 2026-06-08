# confluence whiteboard

Whiteboard reads

## Synopsis

Cloud Whiteboard read operations.

These commands use documented Cloud v2 whiteboards endpoints. Create/delete operations
are mutations and remain deferred behind explicit safety gates.

## Examples

confluence whiteboard view 12345 --json
  confluence whiteboard children 12345 --type page --limit 25

## Usage

```text
confluence whiteboard
```

## Commands

- `children` - List direct children of a Cloud whiteboard
- `view` - Show one Cloud whiteboard

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
