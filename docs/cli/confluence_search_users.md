# confluence search users

Search users via CQL user.fullname~

## Synopsis

Search users via /rest/api/search with
cql = 'type=user AND user.fullname ~ "<query>"'.

## Examples

confluence search users "Jane Smith" --json

## Usage

```text
confluence search users <query> [flags]
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
