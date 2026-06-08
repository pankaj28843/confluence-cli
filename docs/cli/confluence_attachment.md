# confluence attachment

Attachments (list, download, upload, replace, delete)

## Synopsis

Attachment operations.

## Examples

confluence attachment list --page 12345
  confluence attachment download --page 12345 --name logo.png --output ./logo.png
  confluence attachment upload --page 12345 --file ./report.pdf
  confluence attachment replace --page 12345 --file ./report.pdf
  confluence attachment delete --page 12345 --name old.png --force

## Usage

```text
confluence attachment
```

## Commands

- `delete` - Delete an attachment by id or by page/name
- `download` - Download an attachment by name
- `list` - List attachments on a page
- `replace` - Create or replace an attachment on a page
- `upload` - Upload an attachment to a page

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
