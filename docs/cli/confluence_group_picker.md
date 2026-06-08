# confluence group picker

Search Cloud groups by picker query

## Synopsis

Search Cloud groups using the documented group picker endpoint.

## Examples

confluence group picker eng
  confluence group picker eng --limit 50 --total-size --json

## Usage

```text
confluence group picker <query> [flags]
```

## Options

```text
      --limit int    Max groups (hard cap 200) (default 25)
      --total-size   Ask Cloud to include total size metadata
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

- [confluence group](confluence_group.md)
