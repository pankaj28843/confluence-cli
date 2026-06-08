# confluence label add

Add one or more labels to a content id

## Synopsis

Add labels. Pass multiple --label flags or one comma-separated list.

## Examples

confluence label add --page 12345 --label needs-review
  confluence label add --page 12345 --label review,shipped,v1

## Usage

```text
confluence label add [flags]
```

## Options

```text
      --label strings   Label name(s); repeatable
      --page string     Content id (required)
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
