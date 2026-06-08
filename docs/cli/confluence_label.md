# confluence label

Content and space labels

## Synopsis

Label operations.

## Examples

confluence label list --page 12345 --json
  confluence label space --space ENG --limit 50
  confluence label search --prefix global --json
  confluence label add --page 12345 --label needs-review,shipped
  confluence label remove --page 12345 --label needs-review

## Usage

```text
confluence label
```

## Commands

- `add` - Add one or more labels to a content id
- `list` - List labels on a content target
- `recent` - List recently used Server/Data Center labels
- `related` - List related Server/Data Center labels
- `remove` - Remove one label from a content id
- `search` - Search the Cloud label catalog
- `space` - List labels used in a space

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
