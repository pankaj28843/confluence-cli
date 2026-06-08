# confluence page delete

Move a page to trash

## Synopsis

Move a page to trash. On Cloud, --draft deletes a draft page; discarded
drafts are permanently deleted by Confluence and are not sent to trash.

## Examples

confluence page delete 12345
  confluence page delete 12345 --force
  confluence page delete 12345 --draft --force --json

## Usage

```text
confluence page delete <id> [flags]
```

## Options

```text
      --draft   Cloud only: delete a draft page
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

- [confluence page](confluence_page.md)
