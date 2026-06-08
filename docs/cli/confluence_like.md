# confluence like

Cloud likes (count, users)

## Synopsis

Cloud like helpers for pages, blog posts, footer comments, and inline comments.

## Examples

confluence like count --page 12345
  confluence like users --blogpost 67890 --limit 50 --json

## Usage

```text
confluence like
```

## Commands

- `count` - Show Cloud like count for one entity
- `users` - List Cloud account ids that liked one entity

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
