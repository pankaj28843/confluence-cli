# confluence like users

List Cloud account ids that liked one entity

## Synopsis

List Cloud account ids that liked one page, blog post, footer comment, or inline comment.

## Examples

confluence like users --page 12345
  confluence like users --inline-comment 67890 --limit 50 --json

## Usage

```text
confluence like users [flags]
```

## Options

```text
      --blogpost string         Blog post id
      --footer-comment string   Footer comment id
      --inline-comment string   Inline comment id
      --limit int               Max users (hard cap 200) (default 25)
      --page string             Page id
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

- [confluence like](confluence_like.md)
