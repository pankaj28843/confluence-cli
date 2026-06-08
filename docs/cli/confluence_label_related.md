# confluence label related

List related Server/Data Center labels

## Synopsis

List related Server/Data Center labels.

## Examples

confluence label related --label incident --json
  confluence label related --space ENG --label incident --limit 100

## Usage

```text
confluence label related [flags]
```

## Options

```text
      --label string   Label name (required)
      --limit int      Max labels (hard cap 200) (default 25)
      --space string   Server/Data Center space key
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
