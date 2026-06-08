# confluence like count

Show Cloud like count for one entity

## Synopsis

Show the Cloud like count for one page, blog post, footer comment, or inline comment.

## Examples

confluence like count --page 12345
  confluence like count --footer-comment 67890 --json

## Usage

```text
confluence like count [flags]
```

## Options

```text
      --blogpost string         Blog post id
      --footer-comment string   Footer comment id
      --inline-comment string   Inline comment id
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
