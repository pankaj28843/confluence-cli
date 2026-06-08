# confluence folder

Folder reads

## Synopsis

Cloud Folder read operations.

These commands use documented Cloud v2 folders endpoints. Create/delete operations
are mutations and remain deferred behind explicit safety gates.

## Examples

confluence folder view 12345 --json
  confluence folder children 12345 --type page --limit 25

## Usage

```text
confluence folder
```

## Commands

- `children` - List direct children of a Cloud folder
- `view` - Show one Cloud folder

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
