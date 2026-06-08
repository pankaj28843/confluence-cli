# confluence search spaces

Search spaces via CQL type=space

## Synopsis

Search spaces via /rest/api/search with
cql = 'type=space AND (title ~ "Q" OR text ~ "Q")'.

## Examples

confluence search spaces "engineering" --json

## Usage

```text
confluence search spaces <query> [flags]
```

## Options

```text
      --limit int   Max results (hard cap 200) (default 25)
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
