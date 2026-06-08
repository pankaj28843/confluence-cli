# confluence operation list

List permitted operations for one entity

## Synopsis

List permitted operations for one Confluence entity.

## Examples

confluence operation list --page 12345
  confluence operation list --blogpost 67890 --json
  confluence operation list --space ENG

## Usage

```text
confluence operation list [flags]
```

## Options

```text
      --attachment string       Attachment id
      --blogpost string         Blog post id
      --custom-content string   Custom content id
      --database string         Database id
      --embed string            Smart Link embed id
      --folder string           Folder id
      --footer-comment string   Footer comment id
      --inline-comment string   Inline comment id
      --page string             Page id
      --space string            Cloud space key or id
      --whiteboard string       Whiteboard id
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

- [confluence operation](confluence_operation.md)
