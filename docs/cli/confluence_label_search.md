# confluence label search

Search the Cloud label catalog

## Synopsis

Search the Cloud label catalog.

## Examples

confluence label search --json
  confluence label search --prefix global --limit 100
  confluence label search --label-id 123 --json

## Usage

```text
confluence label search [flags]
```

## Options

```text
      --label-id strings   Cloud label id filter; repeatable
      --limit int          Max labels (hard cap 200) (default 25)
      --prefix strings     Cloud label prefix filter; repeatable
      --sort string        Cloud sort expression
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
