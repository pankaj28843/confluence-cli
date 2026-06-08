# confluence group list

List groups

## Synopsis

List groups.

## Examples

confluence group list
  confluence group list --json --limit 200
  confluence group list --access-type admin --json

## Usage

```text
confluence group list [flags]
```

## Options

```text
      --access-type string   Cloud access type filter: user, admin, or site-admin
      --limit int            Max groups (hard cap 200) (default 25)
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
