# confluence page publish

Upload attachments, then update page body

## Synopsis

Publish a page body and any referenced attachments in order. Attachments are
created or updated by filename before the page body is PUT with an incremented
version.

## Examples

confluence page publish 12345 --body-file page.html --attach hld.png
  confluence page publish 12345 --body-file page.html --attach hld.png --attach flow.png
  confluence page publish 12345 --title "Runbook" --body-format storage --body-file page.html

## Usage

```text
confluence page publish <id> [flags]
```

## Options

```text
      --attach stringArray   Attachment file to create or update; repeatable
      --body-file string     Path to body file, or '-' for stdin (required)
      --body-format string   Body format: storage | wiki | view (default "storage")
      --title string         New title (keeps existing if omitted)
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
