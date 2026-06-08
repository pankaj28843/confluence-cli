# confluence search attachments

Search attachments via CQL type=attachment

## Synopsis

Search attachments via /rest/api/content/search with
cql = 'type=attachment AND (title ~ "Q" OR text ~ "Q")'.

## Examples

confluence search attachments "report.pdf" --json
  confluence search attachments "logo" --space ENG

## Usage

```text
confluence search attachments <query> [flags]
```

## Options

```text
      --limit int      Max results (hard cap 200) (default 25)
      --space string   Optional space key filter
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

- [confluence search](confluence_search.md)
