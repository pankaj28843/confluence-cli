# confluence search content

Search pages (type=page by default) via CQL text match

## Synopsis

Search pages via /rest/api/content/search with
cql = 'type=page AND text ~ "<query>"' (plus optional --space).

## Examples

confluence search content "release" --limit 10
  confluence search content "deploy" --space ENG --json

## Usage

```text
confluence search content <query> [flags]
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
