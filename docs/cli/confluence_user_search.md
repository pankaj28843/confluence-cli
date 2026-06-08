# confluence user search

Search users by full name (CQL user.fullname~)

## Synopsis

Search users using user.fullname~"<query>" CQL.

## Examples

confluence user search "Jane Smith"
  confluence user search "Jane" --limit 50 --json

## Usage

```text
confluence user search <query> [flags]
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

- [confluence user](confluence_user.md)
