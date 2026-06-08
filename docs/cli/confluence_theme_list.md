# confluence theme list

List Cloud themes

## Synopsis

List available Confluence Cloud themes.

## Examples

confluence theme list
  confluence theme list --start 25 --limit 25 --json

## Usage

```text
confluence theme list [flags]
```

## Options

```text
      --limit int   Max themes (hard cap 200) (default 25)
      --start int   Zero-based result offset
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

- [confluence theme](confluence_theme.md)
