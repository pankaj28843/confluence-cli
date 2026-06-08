# confluence blogpost delete

Move a blog post to trash

## Synopsis

Move a blog post to trash. On Cloud, --draft deletes a draft blog post;
discarded drafts are permanently deleted by Confluence and are not sent to trash.

## Examples

confluence blogpost delete 12345
  confluence blogpost delete 12345 --force
  confluence blogpost delete 12345 --draft --force --json

## Usage

```text
confluence blogpost delete <id> [flags]
```

## Options

```text
      --draft   Cloud only: delete a draft blog post
      --force   Do not prompt for confirmation
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

- [confluence blogpost](confluence_blogpost.md)
