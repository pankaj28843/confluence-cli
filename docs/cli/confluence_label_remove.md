# confluence label remove

Remove one label from a content id

## Synopsis

Remove a label.

## Examples

confluence label remove --page 12345 --label needs-review

## Usage

```text
confluence label remove [flags]
```

## Options

```text
      --label string   Label name to remove (required)
      --page string    Content id (required)
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

- [confluence label](confluence_label.md)
