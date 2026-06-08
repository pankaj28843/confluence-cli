# confluence attachment delete

Delete an attachment by id or by page/name

## Synopsis

Delete an attachment content entity. Pass --id directly, or pass --page and
--name to resolve the attachment id first. A confirmation prompt is shown unless
--force is supplied.

## Examples

confluence attachment delete --id 1884909332 --force
  confluence attachment delete --page 12345 --name old-diagram.png
  confluence attachment delete --page 12345 --name old-diagram.png --force --json

## Usage

```text
confluence attachment delete [flags]
```

## Options

```text
      --force         Delete without confirmation
      --id string     Attachment content id
      --name string   Attachment filename used with --page
      --page string   Page id used with --name
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

- [confluence attachment](confluence_attachment.md)
