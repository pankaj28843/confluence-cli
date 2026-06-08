# confluence analytics

Cloud content analytics reads

## Synopsis

Cloud content analytics read operations.

Analytics view and viewer counts are documented in Confluence Cloud REST API v1.
Server/Data Center does not expose the same Analytics REST group in the current
official REST reference, so these typed commands are Cloud-only.

## Examples

confluence analytics views 12345 --json
  confluence analytics viewers 12345 --from-date YYYY-MM-DDTHH:MM:SS.sssZ

## Usage

```text
confluence analytics
```

## Commands

- `viewers` - Show Cloud distinct viewer count
- `views` - Show Cloud content view count

## Inherited Options

```text
      --debug             Log HTTP requests to stderr (Authorization header redacted)
      --jq string         Filter JSON output through a jq expression (requires jq on PATH)
      --json              Output as JSON (machine-readable)
      --template string   Render JSON output through a Go text/template
      --timing            Show execution time on stderr
```

## See Also

- [confluence](confluence.md)
