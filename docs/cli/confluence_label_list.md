# confluence label list

List labels on a content id

## Synopsis

List labels.

## Examples

confluence label list --page 12345 --json
  confluence label list --page 12345 --limit 100

## Usage

```text
confluence label list [flags]
```

## Options

```text
      --limit int     Max labels (hard cap 200) (default 25)
      --page string   Content id (required)
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
