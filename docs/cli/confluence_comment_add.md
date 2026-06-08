# confluence comment add

Add a footer comment

## Synopsis

Add a footer comment to a page or blog post. On Confluence Cloud, --parent
creates a reply to an existing footer comment.

## Examples

confluence comment add --page 12345 --body "<p>Looks good.</p>"
  confluence comment add --blogpost 2001 --body-file comment.html
  echo "<p>Reply</p>" | confluence comment add --parent 998877 --body-file - --json

## Usage

```text
confluence comment add [flags]
```

## Options

```text
      --blogpost string      Blog post id to comment on
      --body string          Inline body string
      --body-file string     Path to body file, or '-' for stdin
      --body-format string   Body format: storage (default "storage")
      --page string          Page id to comment on
      --parent string        Cloud only: parent footer comment id for a reply
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

- [confluence comment](confluence_comment.md)
