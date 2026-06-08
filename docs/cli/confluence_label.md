# confluence label

Content labels (list, add, remove)

## Synopsis

Label operations.

## Examples

confluence label list --page 12345
  confluence label add --page 12345 --label needs-review,shipped
  confluence label remove --page 12345 --label needs-review

## Usage

```text
confluence label
```

## Commands

- `add` - Add one or more labels to a content id
- `list` - List labels on a content id
- `remove` - Remove one label from a content id

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
