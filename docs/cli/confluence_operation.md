# confluence operation

Permitted operations (list)

## Synopsis

Permitted operation helpers.

Cloud supports pages, blog posts, attachments, spaces, comments, and newer
content-tree entities. Server/Data Center supports content ids through the
documented operations expansion.

## Examples

confluence operation list --page 12345
  confluence operation list --space ENG --json

## Usage

```text
confluence operation
```

## Commands

- `list` - List permitted operations for one entity

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
