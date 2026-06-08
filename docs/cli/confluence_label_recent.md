# confluence label recent

List recently used Server/Data Center labels

## Synopsis

List recently used Server/Data Center labels.

## Examples

confluence label recent --json
  confluence label recent --limit 100

## Usage

```text
confluence label recent [flags]
```

## Options

```text
      --limit int   Max labels (hard cap 200) (default 25)
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
