# confluence label list

List labels on a content target

## Synopsis

List labels.

## Examples

confluence label list --page 12345 --json
  confluence label list --blogpost 67890 --prefix global
  confluence label list --attachment att123 --limit 100

## Usage

```text
confluence label list [flags]
```

## Options

```text
      --attachment string       Attachment id
      --blogpost string         Blog post id
      --custom-content string   Custom content id
      --limit int               Max labels (hard cap 200) (default 25)
      --page string             Page id
      --prefix string           Label prefix filter
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
