# confluence label space

List labels used in a space

## Synopsis

List labels used in a space.

## Examples

confluence label space --space ENG --json
  confluence label space --space ENG --prefix global --limit 100
  confluence label space --space ENG --scope space --json

## Usage

```text
confluence label space [flags]
```

## Options

```text
      --limit int       Max labels (hard cap 200) (default 25)
      --prefix string   Label prefix filter
      --scope string    Cloud only: content or space (default "content")
      --space string    Space key or Cloud space id (required)
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
