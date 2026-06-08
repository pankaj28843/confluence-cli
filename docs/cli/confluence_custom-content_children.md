# confluence custom-content children

List child Cloud custom content

## Synopsis

List child custom content under one Cloud custom-content id.

## Examples

confluence custom-content children 12345
  confluence custom-content children 12345 --type ac:example --json
  confluence custom-content children 12345 --sort title --limit 25

## Usage

```text
confluence custom-content children <id> [flags]
```

## Options

```text
      --limit int      Max children (hard cap 200) (default 25)
      --sort string    Cloud sort expression
      --type strings   Client-side custom content type filter; repeatable or comma-separated
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

- [confluence custom-content](confluence_custom-content.md)
