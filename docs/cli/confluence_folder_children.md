# confluence folder children

List direct children of a Cloud folder

## Synopsis

List direct children of one Cloud folder.

The Cloud v2 endpoint returns minimal child rows for databases, Smart Link
embeds, folders, pages, and whiteboards. Use the matching view command for
full details.

## Examples

confluence folder children 12345
  confluence folder children 12345 --type page --type database --json
  confluence folder children 12345 --sort position --limit 100

## Usage

```text
confluence folder children <id> [flags]
```

## Options

```text
      --limit int      Max children (hard cap 200) (default 50)
      --sort string    Cloud sort expression
      --type strings   Content type filter; repeatable or comma-separated
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

- [confluence folder](confluence_folder.md)
