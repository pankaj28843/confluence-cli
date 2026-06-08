# confluence attachment list

List attachments on a page

## Synopsis

List attachments.

## Examples

confluence attachment list --page 12345
  confluence attachment list --page 12345 --name hld.png
  confluence attachment list --page 12345 --json --limit 100

## Usage

```text
confluence attachment list [flags]
```

## Options

```text
      --limit int     Max attachments (hard cap 200) (default 50)
      --name string   Exact attachment file name filter
      --page string   Content id (required)
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
