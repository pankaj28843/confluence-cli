# confluence body

Content body utilities

## Synopsis

Content body utility operations.

## Examples

confluence body convert --from storage --to view --value '<p>Hello</p>'
  confluence body convert --to export_view --value @body-storage.xml --json

## Usage

```text
confluence body
```

## Commands

- `convert` - Convert a content body representation

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
