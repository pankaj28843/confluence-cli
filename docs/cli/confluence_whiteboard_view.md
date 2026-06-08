# confluence whiteboard view

Show one Cloud whiteboard

## Synopsis

Show one Cloud whiteboard by id.

## Examples

confluence whiteboard view 12345
  confluence whiteboard view 12345 --include-operations --include-properties --json

## Usage

```text
confluence whiteboard view <id> [flags]
```

## Options

```text
      --include-collaborators     Include collaborators in the Cloud response
      --include-direct-children   Include direct children in the Cloud response
      --include-operations        Include permitted operations in the Cloud response
      --include-properties        Include properties in the Cloud response
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

- [confluence whiteboard](confluence_whiteboard.md)
