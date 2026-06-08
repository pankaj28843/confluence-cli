# confluence analytics views

Show Cloud content view count

## Synopsis

Show the total number of Cloud views for one content item.

## Examples

confluence analytics views 12345
  confluence analytics views 12345 --from-date YYYY-MM-DDTHH:MM:SS.sssZ --json

## Usage

```text
confluence analytics views <content-id> [flags]
```

## Options

```text
      --from-date string   Cloud analytics start date, e.g. YYYY-MM-DDTHH:MM:SS.sssZ
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

- [confluence analytics](confluence_analytics.md)
